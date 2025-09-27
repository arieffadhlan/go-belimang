package entities

type User struct {
	Id       string `db:"id"`
	Email    string `db:"email"`
	Username string `db:"username"`
	Password string `db:"password"`
	IsAdmin  bool   `db:"is_admin"`
}
