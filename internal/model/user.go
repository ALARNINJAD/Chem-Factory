package model

type User struct {
	ID       int    `db:"id" json:"id"`
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"password"`
	Balance  int    `db:"balance" json:"balace"`
	XP       int    `db:"xp" json:"xp"`
	Level    int    `db:"level" json:"level"`
}
