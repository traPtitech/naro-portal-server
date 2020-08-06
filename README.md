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
### users
| Field       | Type         | NULL | Key | Default | Extra | 説明                           |
|-------------|--------------|------|-----|---------|-------|--------------------------------|
| id          | varchar(36)  | NO   | PRI |         |       | userのid,tweets.authorと同じ   |
| hashed_pass | varchar(256) | NO   |     |         |       | ハッシュ化されたパスワード     |

### tweets
| Field      | Type           | NULL   | Key   | Default   | Extra   | 説明               |
| ---------- | -------------- | ------ | ----- | --------- | ------- | ------------------ |
| id         | char(36)       | NO     | PRI   |           |         | tweetのuuid        |
| user_id    | varchar(32)    | NO     |       |           |         | ツイートした人のid |
| tweet_body | varchar(256)   | NO     |       |           |         | 本文               |
| created_at | datetime       | NO     |       |           |         | ツイート時刻       |
