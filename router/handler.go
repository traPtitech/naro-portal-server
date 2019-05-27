package router

//ハンドラ
import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/WistreHosshii/naro-portal-server/model/users"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

var (
	db *sqlx.DB
)

func Pong(c echo.Context) error {
	fmt.Println(c)
	return c.String(http.StatusOK, "pong")
}

func PostSignUpHandler(c echo.Context) error {
	req := users.LoginReqestBody{}
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
	var id string
	//idの生成。idが被った場合5回まで作り直す
	for i := 0; i < 5; i++ {
		id = generateId(req.UserName)
		err = db.Get(&count, "SELECT COUNT(*) FROM users WHERE id=?", id)
		if err != nil {

			return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
		}
		//idが被らなかったらok
		if count == 0 {
			_, err = db.Exec("INSERT INTO users (Username, HashedPass, ID) VALUES (?, ?, ?)", req.UserName, hashedPass, id)
			if err != nil {
				return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
			}
			return c.NoContent(http.StatusCreated)
		}
	}
	return c.String(http.StatusInternalServerError,"idが生成できません")
}

func generateId(UserName string) string {
	rand.Seed(time.Now().UnixNano())
	var r = int(rand.Float64() * 100000)
	var id = UserName + strconv.Itoa(r)
	return id
}
