package users

type User struct {
	UserName   string `json:"username,omitempty" db:user_name`
	HashadPass string `json:"hashedpass,omitempty" db:hashed_pass`
	ID         string `json:id,omitempty" db:id`
}

type LoginReqestBody struct {
	UserName string `json:"username,omitempty" db:user_name`
	Password string `json:"password,omitempty" db:password`
}
//CREATE TABLE `users` (`user_name` VARCHAR(20) NOT NULL, `hashed_pass` VARCHAR(200) NOT NULL, `id` VARCHAR(100) NOT NULL, PRIMARY KEY(`id`)) ENGINE = InnoDB;