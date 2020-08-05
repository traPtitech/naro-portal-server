# twitter-clone
Webエンジニアになろう講習会課題のサーバリポジトリ

## EndPoint
### 非ログイン時
* `POST /register` ユーザ登録
* `POST /login` ログイン
* `GET /tweet` ツイート取得(30件)

### ログイン時
* `GET /timeline` ツイート取得(30件)
* `POST /timeline` ツイート投稿

## DB schema
### tweets
| Field      | Type           | NULL   | Key   | Default   | Extra   | 説明             |
| ---------- | -------------- | ------ | ----- | --------- | ------- | ---------------- |
| id         | char(36)       | NO     | PRI   |           |         | tweetのuuid      |
| tweet_body | varchar(256)   | NO     |       |           |         | 本文             |
| author     | varchar(32)    | NO     |       |           |         | ツイートした人   |
| created_at | datetime       | NO     |       |           |         | ツイート時刻     |
