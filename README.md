# kuragate-server

Web エンジニアになろう講習会課題のポータルリポジトリ(Server)

## APIs

| URL                     | Method | Description                           | Implemented |
| ----------------------- | ------ | ------------------------------------- | ----------- |
| /ping                   | GET    | Return "pong"                         | o           |
| /isvalidid/:id          | GET    | Return true if given id is valid      | o           |
| /signup                 | POST   | Signup                                | o           |
| /login                  | POST   | Login                                 | o           |
| /whoami                 | GET    | Return profiles of logged in user     | o           |
| /logout                 | GET    | Logout                                | o           |
| /messages               | POST   | Post message                          | o           |
| /messages               | GET    | Get all messages                      | o           |
| /messages/:id           | GET    | Get single message which has given id | o           |
| /messages/:id/fav       | PUT    | Fav a message                         | o           |
| /messages/:id/fav       | DELETE | Unfav a message                       | o           |
| /messages/:id/fav       | GET    | Get users who fav a message           | x           |
| /profiles/:id           | GET    | Return profile of a user              | x           |
| /profiles/:id/messages  | GET    | Return messages posted by a user      | x           |
| /profiles/:id/following | GET    | Return users followed by a user       | x           |
| /profiles/:id/followed  | GET    | Return users who are following a user | x           |
| /profiles/:id/followed  | PUT    | follow user                           | x           |
| /profiles/:id/followed  | DELETE | unfollow user                         | x           |

GET /ping

    e.GET("/ping", func(c echo.Context) error {
    	return c.String(http.StatusOK, "pong")
    })

    e.GET("/isvalidid/:reqID", auths.GetIsValidIDHandler)
    e.POST("/signup", auths.PostSignUpHandler)
    e.POST("/login", auths.PostLoginHandler)

    withLogin := e.Group("")
    withLogin.Use(auths.CheckLogin)

    withLogin.GET("/whoami", auths.GetWhoAmIHandler)
    withLogin.GET("/logout", auths.GetLogoutHandler)

    withLogin.POST("/messages", messages.PostMessageHandler)
    withLogin.GET("/messages", messages.GetMassagesHandler)
    withLogin.GET("/messages/:id", messages.GetSingleMassageHandler)
    withLogin.PUT("/messages/:id/fav", messages.PutMessageFavHandler)
    withLogin.DELETE("/messages/:id/fav", messages.DeleteMessageFavHandler)
