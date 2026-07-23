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

func (user User) New(username, password string) (User, error) {
	if username == "" {
		return User{}, errors.New("username is invalid")
	}
	if password == "" {
		return User{}, errors.New("password is invalid")
	}
	return User{
		ID:       0,
		Username: username,
		Password: password,
		Balance:  500,
		XP:       0,
		Level:    0,
	}, nil
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
