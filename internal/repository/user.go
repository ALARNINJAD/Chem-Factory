package repository

type user struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Balance  int    `json:"balance"`
	XP       int    `json:"xp"`
	Level    int    `json:"level"`
}

func createUserTable(r *repositoryManager) {
	query := `
	CREATE TABLE IF NOT EXISTS user (
		id INTEGER NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		balance INTEGER DEFAULT 100 CHECK("balance" >= 0),
		xp INTEGER DEFAULT 0 CHECK("xp" >= 0),
		level INTEGER DEFAULT 0 CHECK("level" >= 0)
	)`
	if _, err := r.db.Exec(query); err != nil {
		panic(err.Error())
	}
}

func (r *repositoryManager) ExportUserByUsername(username string) (*user, error) {
	u := user{Username: username}
	err := r.db.QueryRow("SELECT id,password,balance,xp,level FROM user WHERE username = ?", u.Username).Scan(
		&u.ID, &u.Password, &u.Balance, &u.XP, &u.Level)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repositoryManager) ExportUserByID(id int) (*user, error) {
	u := user{ID: id}
	err := r.db.QueryRow("SELECT username,password,balance,xp,level FROM user WHERE id = ?", u.ID).Scan(
		&u.Username, &u.Password, &u.Balance, &u.XP, &u.Level)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repositoryManager) ExportPasswordByID(id int) (string, error) {
	var savedPassword string
	err := r.db.QueryRow("SELECT password FROM user WHERE id = ?", id).Scan(&savedPassword)
	return savedPassword, err
}

func (r *repositoryManager) ExportPasswordByUsername(username string) (string, error) {
	var savedPassword string
	err := r.db.QueryRow("SELECT password FROM user WHERE username = ?", username).Scan(&savedPassword)
	return savedPassword, err
}

func (r *repositoryManager) ExportIDbyUsername(username string) (int, error) {
	var id int
	err := r.db.QueryRow("SELECT id FROM user WHERE username = ?", username).Scan(&id)
	return id, err
}

func (r *repositoryManager) ExportUsernameByID(id int) (string, error) {
	var username string
	err := r.db.QueryRow("SELECT username FROM user WHERE id = ?", id).Scan(&username)
	return username, err
}

func (r *repositoryManager) SaveNewUser(username string, password string) error {
	_, err := r.db.Exec("INSERT INTO user(username, password) VALUES (?, ?)", username, password)
	return err
}

func (r *repositoryManager) IncreaseBalance(username string, amount int) error {

	_, err := r.db.Exec("UPDATE user SET balance = balance + ? WHERE username = ?", amount, username)
	if err != nil {
		return err
	}

	return nil
}

func (r *repositoryManager) ReduceBalance(username string, amount int) error {

	_, err := r.db.Exec("UPDATE user SET balance = balance - ? WHERE username = ?", amount, username)
	if err != nil {
		return err
	}

	return nil
}
