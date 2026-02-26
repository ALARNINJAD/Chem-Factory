package repository

import "chem-factory/internal/model"

func Export(user *model.User) error {
	if user.ID > 0 {
		if user.Username != "" {
			return db.QueryRow("SELECT password,balance,xp,level FROM user WHERE id = ?", user.ID).Scan(
				&user.Password, &user.Balance, &user.XP, &user.Level)
		}
		return db.QueryRow("SELECT username,password,balance,xp,level FROM user WHERE id = ?", user.ID).Scan(
			&user.Username, &user.Password, &user.Balance, &user.XP, &user.Level)
	}
	return db.QueryRow("SELECT id,password,balance,xp,level FROM user WHERE username = ?", user.Username).Scan(
		&user.ID, &user.Password, &user.Balance, &user.XP, &user.Level)
}

func Save(user *model.User) error {
	_, err := db.Exec("INSERT INTO user(username, password) VALUES (?, ?)", user.Username, user.Password)
	return err
}
