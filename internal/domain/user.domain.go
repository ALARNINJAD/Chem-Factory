package domain

import "errors"

type User struct {
	ID       uint   `json:"id,omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
	Balance  int    `json:"balance"`
	XP       int    `json:"xp"`
	Level    int    `json:"level"`
}

func (user *User) New(username, password string) error {
	if username == "" {
		return errors.New("username is invalid")
	}
	if password == "" {
		return errors.New("password is invalid")
	}
	user.Username = username
	user.Password = password
	user.Balance = 0
	user.XP = 0
	user.Level = 1
	return nil
}

func (user *User) IncreaseBalance(amount int) error {
	if amount < 0 {
		return errors.New("increase amount must be positive")
	}
	user.Balance += amount
	return nil
}

func (user *User) ReduceBalance(amount int) error {
	if amount < 0 {
		return errors.New("reduce amount must be positive")
	}
	if user.Balance < amount {
		return errors.New("insufficient balance")
	}
	user.Balance -= amount
	return nil
}

func (user *User) GetXP(xp int) error {
	if xp < 0 {
		return errors.New("increase amount must be positive")
	}
	user.XP += xp
	for user.XP >= user.Level*1000 {
		user.XP -= user.Level * 1000
		user.Level++
	}
	return nil
}
