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
