package router

//ハンドラ
import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/WistreHosshii/naro-portal-server/model"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/oklog/ulid"

	//"strconv"

	//"github.com/WistreHosshii/naro-portal-server/model"
	"github.com/WistreHosshii/naro-portal-server/model/mystruct"

	_ "github.com/go-sql-driver/mysql"

	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

func Pong(c echo.Context) error {
	fmt.Println(c)
	return c.String(http.StatusOK, "pong")
}

func PostSignUpHandler(c echo.Context) error {
	req := mystruct.LoginReqestBody{}
	c.Bind(&req)

	if req.Password == "" || req.UserName == "" { //パスワードとか名前が空？
		return c.String(http.StatusBadRequest, "項目が空です")
	} else if len(req.Password) < 5 {
		//パスワード短くないか検証
		return c.String(http.StatusBadRequest, "パスワードが短すぎます")
	}

	//上のパスワードとかの処理はクライアント側でやったほうが良い気がする

	//ハッシュ化
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("bcrypt generate error: %v", err))
	}

	//同じユーザー名が何人いるか調べる
	var count int
	count, err = model.GetUserCount(req.UserName)

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}

	if count > 0 {
		return c.String(http.StatusConflict, "ユーザーが既に存在しています")
	}

	id := ExampleULID

	//idの生成。idが被った場合5回まで作り直す
	/*for i := 0; i < 5; i++ {
		id = generateID(req.UserName)
		err = Db.Get(&count, "SELECT COUNT(*) FROM users WHERE id=?", id)
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
		}
		//idが被らなかったらok
		if count == 0 {
			_, err = Db.Exec("INSERT INTO users (user_name, hashed_pass, id) VALUES (?, ?, ?)", req.UserName, hashedPass, id)
			if err != nil {
				return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
			}
			return c.NoContent(http.StatusCreated)
		}
	}
	return c.String(http.StatusInternalServerError, "idが生成できません")
	*/
	model.ExecUserInfo(req.UserName, hashedPass, id)
	return c.NoContent(http.StatusCreated)
}

/*func generateID(UserName string) string {
	rand.Seed(time.Now().UnixNano())
	var r = int(rand.Float64() * 1000000000)
	var id = UserName + strconv.Itoa(r)
	return id
}
*/
func ExampleULID() func() string {
	t := time.Unix(1000000, 0)
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	//fmt.Println(ulid.MustNew(ulid.Timestamp(t), entropy))
	// Output: 0000XSNJG0MQJHBF4QX1EFD6Y3
	_id := ulid.MustNew(ulid.Timestamp(t), entropy).String
	return _id

}

func PostLoginHandler(c echo.Context) error {
	req := mystruct.LoginReqestBody{}
	c.Bind(&req)

	//user := users.User{}
	//err := Db.Get(&user, "SELECT FROM users WHERE user_name=?",req.UserName)
	user, err := model.GetUserName(req.UserName)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPass), []byte(req.Password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return c.NoContent(http.StatusForbidden)
		} else {
			return c.NoContent(http.StatusInternalServerError)
		}
	}

	sess, err := session.Get("sessions", c)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "something wrong in getting session")
	}

	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	sess.Values["userName"] = req.UserName
	sess.Save(c.Request(), c.Response())

	return c.NoContent(http.StatusOK)
}
