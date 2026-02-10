package domain

type User struct {
	UserID   int    `db:"user_id" json:"user_id"`
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"password"`
	Balance  int    `db:"balance" json:"balance"`
	XP       int    `db:"xp" json:"xp"`
	Level    int    `db:"level" json:"level"`
}
