package repository

type inventory struct {
	id         int `db:"id"`
	userID     int `db:"user_id"`
	materialID int `db:"material_id"`
	number     int `db:"number"`
}
