package main

import (
	"fmt"
	"log"
	"naro-server/go/src/naro-server/packages/handler"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/srinathgs/mysqlstore"
	"golang.org/x/crypto/acme/autocert"
)

var (
	db *sqlx.DB
)

//var sessionManager *session.Manager;

func main() {
	// if len(os.Args) != 2 {
	// 	log.Println("街の名前を一つ渡してください")
	// 	return
	// }
	//cityName := os.Args[1]
	//cityName := "traP"
	_db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOSTNAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE")))

	if err != nil {
		log.Fatalf("Cannot Connect to Database: %s", err)
	}
	fmt.Println("Connected!")

	db = _db

	store, err := mysqlstore.NewMySQLStoreFromConnection(db.DB, "sessions", "/", 60*60*24*14, []byte("secret-token"))
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(session.Middleware(store))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})
	e.POST("/login", handler.PostLoginHandler(db))
	e.POST("/signup", handler.PostSignUpHandler(db))
	e.GET("/logout", handler.LogoutHandler(db))
	e.GET("/postList", handler.GetPostListHandler(db))
	e.GET("/post/:postID", handler.GetPostHandler(db))

	withLogin := e.Group("")
	withLogin.Use(handler.CheckLogin)

	withLogin.GET("/whoami", handler.GetWhoAmIHandler(db))

	withLogin.POST(":traQID/tasks", handler.AddTaskHandler(db))
	withLogin.GET(":traQID/tasks", handler.GetTaskHandler(db))
	withLogin.DELETE(":traQID/tasks", handler.DeleteAllTasksHandler(db))
	withLogin.DELETE(":traQID/tasks/:ID", handler.DeleteTaskHandler(db))
	withLogin.PUT(":traQID/tasks/:ID", handler.PutTaskHandler(db))

	withLogin.GET("/cities/:cityName", handler.GetCityInfoHandler(db))
	withLogin.GET("/countries", handler.GetCountryListHandler(db))
	withLogin.GET("/countries/:countryCode/cities", handler.GetCityListInCountryHandler(db))
	withLogin.GET("/getUserName", handler.GetUserNameHandler())
	withLogin.GET("/delete/:cityID", handler.DeleteCityInfoHandler(db))
	withLogin.POST("/addCity", handler.AddCityHandler(db))
	withLogin.POST("/post", handler.AddPostHandler(db))
	e.Logger.Fatal(e.Start(":10101"))

	//cityState := `INSERT INTO city(ID, Name, CountryCode, District, Population) VALUES (5001,'traP','JPN','titech',10000000)`
	//db.Exec(cityState)

	// var city City
	// if err := db.Get(&city, "SELECT * FROM city WHERE Name='"+cityName+"'"); errors.Is(err, sql.ErrNoRows) {
	// 	log.Printf("no such city Name = %s", cityName)
	// 	return
	// } else if err != nil {
	// 	log.Fatalf("DB Error: %s", err)
	// }
	// countryCode := city.CountryCode
	// var country Country
	// if err := db.Get(&country, "SELECT * FROM country WHERE Code='"+countryCode+"'"); errors.Is(err, sql.ErrNoRows) {
	// 	log.Printf("no such country code = %s", countryCode)
	// 	return
	// } else if err != nil {
	// 	log.Fatalf("DB Error: %s", err)
	// }

	// fmt.Printf("%sの人口は%d人です\n", cityName, city.Population)
	// fmt.Printf("%sの人口は%sの%.4f%%人です\n", cityName, country.Name, (float64(city.Population)/float64(country.Population))*100)
}
