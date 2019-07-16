package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/srinathgs/mysqlstore"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type userRequestBody struct {
	Name        string    `json:"name,omitempty"`
	Password    string    `json:"password,omitempty"`
}

type tweetRequestBody struct {
	Text   string `json:"text,omitempty"`
	UserID int    `json:"user_id,omitempty"`
	ImgIDs []int  `json:"img_ids,omitempty"`
}

var (
	db *sqlx.DB
)

func main() {
	_db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOSTNAME"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE")))
	if err != nil {
		log.Fatalf("Cannot Connect to Database: %s", err)
	} //DataBaseに接続してる（エラーならその内容を表示）
	db = _db //func mainの外でも使えるように外で定義したdbに_dbを代入

	store, err := mysqlstore.NewMySQLStoreFromConnection(db.DB, "sessions", "/", 60*60*24*14, []byte("secret-token"))
	if err != nil {
		panic(err)
	} //store, errを上手いことしてくれる　panic=実行をその時点で終了する

	e := echo.New()            //echoのインスタンス（echo＝サーバーに関するリクエストとかレスポンスとかの情報諸々を処理してくれるライブラリ）
	e.Use(middleware.Logger()) //通行者のlogをとる
	//e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:1234"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	}))
	e.Use(session.Middleware(store)) //通行証の正当性を確認したのち、echo.Contextにその情報を追加
	e.Use(middleware.Static("."))

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	}) //ピンポン

	e.Static("/", "public")
	//file upload
	e.POST("/img", postImgHandler)
	e.POST("/tweet", postTweetHandler)

	//.GET("/search", getTweetsHandler)

	e.POST("/tweets/:tweet_id/fav", postFavHandler)

	e.POST("/signup", postSignUpHandler)
	e.POST("/login", postLoginHandler) //こいつらはボタンで分かれてる（ちなみにUseは上から順）

	e.Logger.Fatal(e.Start(":80"))
}

func postSignUpHandler(c echo.Context) error {
	req := userRequestBody{}
	c.Bind(&req) //対応するリクエストのKeyの値を構造体にうまくあてはめてくれる

	// もう少し真面目にバリデーションするべき（場合分けがなんとなくガバそう）
	if req.Password == "" || req.Name == "" {
		// エラーは真面目に返すべき
		return c.String(http.StatusBadRequest, "項目が空です")
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost) //bcryptでpasswordをハッシュ化されたパスワードを生成
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("bcrypt generate error: %v", err))
	}

	// 作ろうとしているUserNameが既存のものと重複していないかチェック
	var count int

	err = db.Get(&count, "SELECT COUNT(*) FROM user WHERE name=?", req.Name) //COUNT該当する行数を返す
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}

	if count > 0 {
		return c.String(http.StatusConflict, "ユーザーが既に存在しています")
	}

	//_=返り値が返ってくるけどいらないから明示的に捨てる
	_, err = db.Exec("INSERT INTO user (name, hashed_pass VALUES (?, ?)", req.Name, hashedPass)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}
	return c.NoContent(http.StatusCreated)
}

func postLoginHandler(c echo.Context) error {
	req := userRequestBody{}
	c.Bind(&req)

	type User struct {
		ID          int       `db:"id"`
		Name        string    `db:"name"`
		HashedPass  string    `db:"hashed_pass"`
		CreatedAt   time.Time `db:"created_at"`
		UpdatedAt   time.Time `db:"updated_at"`
	}
	user := User{}
	err := db.Get(&user, "SELECT * FROM user WHERE name=?", req.Name) //リクエストのUserがDBに存在するか問い合わせしていればその情報をUserにその情報を追加
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPass), []byte(req.Password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return c.NoContent(http.StatusForbidden) //Passwordの不一致
		} else {
			return c.NoContent(http.StatusInternalServerError)
		}
	}

	sess, err := session.Get("sessions", c) //session（サーバー側に存在する帳簿）にUserを登録
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "something wrong in getting session")
	}
	sess.Values["userName"] = req.Name
	sess.Save(c.Request(), c.Response())

	return c.NoContent(http.StatusOK)
}

func checkLogin(next echo.HandlerFunc) echo.HandlerFunc { //middlewareと呼ばれるRequestとHandler関数をつなぐもの
	return func(c echo.Context) error {
		sess, err := session.Get("sessions", c) //sessionの取得（userNameがその中に存在しているかどうか？）
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusInternalServerError, "something wrong in getting session")
		}

		if sess.Values["userName"] == nil {
			return c.String(http.StatusForbidden, "please login")
		}
		c.Set("userName", sess.Values["userName"].(string))

		return next(c)
	}
}

func postImgHandler(c echo.Context) error {
	//-----------
	// Read file
	//-----------

	// Source
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	result, err := db.Exec("INSERT INTO img (path) VALUES (?)", "") //空のデータを挿入
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}
	imgID, err := result.LastInsertId() //result型の中身を見る
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}
	path := "data/" + strconv.Itoa(int(imgID)) + ".jpg" //拡張子どうする？

	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	//img tableの更新
	_, err = db.Exec("UPDATE img SET path=? WHERE img.id=?;", path, imgID)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}

	//imgIDを返す
	type Img struct {
		ID        int       `db:"id"`
		Path      string    `db:"path"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}
	uploadedImg := Img{}
	db.Get(&uploadedImg, "SELECT * FROM img WHERE img.id=?", imgID)
	if uploadedImg.Path == "" {
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, uploadedImg)
}

func postTweetHandler(c echo.Context) error {
	tweet := tweetRequestBody{}
	c.Bind(&tweet)

	//tweet.UserID = 1  (うまくidがbindされない)

	result, err := db.Exec("INSERT INTO tweet (text, user_id) VALUES (?, ?)", tweet.Text, tweet.UserID)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}
	tweetID, err := result.LastInsertId() //result型の中身を見る
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}

	for imgID := range tweet.ImgIDs {
		_, err := db.Exec("INSERT INTO test_files (tweet_id, img_id) VALUES (?, ?)", tweetID, imgID)
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
		}
	}

	return c.NoContent(http.StatusCreated)
}

func postFavHandler(c echo.Context) error {
	userName := c.Get("userName").(string)
	tweetID := c.Param("tweet_id")

	_, err := db.Exec("INSERT INTO fav (tweet_id, user_id) VALUE(?, ?)", tweetID, userName)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}
	return c.NoContent(http.StatusCreated)
}
