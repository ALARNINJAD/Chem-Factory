package sqlite

import (
	database "chem-factory/internal/database/sqlite"
	"chem-factory/internal/domain"
	"chem-factory/pkg/reedam"
	"context"
	"errors"
	"gorm.io/gorm"
)

type User struct {
	ID       uint   `gorm:"primaryKey;autoIncrement;not null;unique"`
	Username string `gorm:"not null;unique"`
	Password string `gorm:"not null"`
	Balance  int    `gorm:"default:0;check:balance >= 0"`
	XP       int    `gorm:"default:0;check:xp >= 0"`
	Level    int    `gorm:"default:1;check:level >= 1"`
}

func (User) TableName() string {
	return "users"
}

func toDomainUser(u User) domain.User {
	return domain.User{
		ID:       u.ID,
		Username: u.Username,
		Password: u.Password,
		Balance:  u.Balance,
		XP:       u.XP,
		Level:    u.Level,
	}
}

func fromDomainUser(u domain.User) User {
	return User{
		ID:       u.ID,
		Username: u.Username,
		Password: u.Password,
		Balance:  u.Balance,
		XP:       u.XP,
		Level:    u.Level,
	}
}

type UserRepo struct{ db *database.Database }

func NewUserRepo(db *database.Database) *UserRepo { return &UserRepo{db: db} }

func (r *UserRepo) FindByUsername(ctx context.Context, username string) (domain.User, error) {
	var u User
	err := r.db.Extract(ctx).Where("username = ?", username).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, reedam.UserNotFound(err)
		}
		return domain.User{}, reedam.Unexpected(err).WithLog(username, database.IsTx(ctx))
	}
	return toDomainUser(u), nil
}

func (r *UserRepo) FindByID(ctx context.Context, id uint) (domain.User, error) {
	var u User
	err := r.db.Extract(ctx).First(&u, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, reedam.UserNotFound(err)
		}
		return domain.User{}, reedam.Unexpected(err).WithLog(id, database.IsTx(ctx))
	}
	return toDomainUser(u), nil
}

func (r *UserRepo) FindPasswordByID(ctx context.Context, id uint) (string, error) {
	var u User
	err := r.db.Extract(ctx).Select("password").First(&u, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", reedam.UserNotFound(err)
		}
		return "", reedam.Unexpected(err).WithLog(id, database.IsTx(ctx))
	}
	return u.Password, nil
}

func (r *UserRepo) FindPasswordByUsername(ctx context.Context, username string) (string, error) {
	var u User
	err := r.db.Extract(ctx).Select("password").Where("username = ?", username).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", reedam.UserNotFound(err)
		}
		return "", reedam.Unexpected(err).WithLog(username, database.IsTx(ctx))
	}
	return u.Password, nil
}

func (r *UserRepo) FindIDByUsername(ctx context.Context, username string) (uint, error) {
	var u User
	err := r.db.Extract(ctx).Select("id").Where("username = ?", username).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, reedam.UserNotFound(err)
		}
		return 0, reedam.Unexpected(err).WithLog(username, database.IsTx(ctx))
	}
	return u.ID, nil
}

func (r *UserRepo) FindUsernameByID(ctx context.Context, id uint) (string, error) {
	var u User
	err := r.db.Extract(ctx).Select("username").First(&u, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", reedam.UserNotFound(err)
		}
		return "", reedam.Unexpected(err).WithLog(id, database.IsTx(ctx))
	}
	return u.Username, nil
}

func (r *UserRepo) Add(ctx context.Context, user domain.User) error {
	u := fromDomainUser(user)
	err := r.db.Extract(ctx).Create(&u).Error
	if err != nil {
		return reedam.Unexpected(err).WithLog(user, database.IsTx(ctx))
	}
	return nil
}

func (r *UserRepo) IncreaseBalanceByID(ctx context.Context, id uint, amount int) error {
	err := r.db.Extract(ctx).Model(&User{}).Where("id = ?", id).Update("balance", gorm.Expr("balance + ?", amount)).Error
	if err != nil {
		return reedam.Unexpected(err).WithLog(id, amount, database.IsTx(ctx))
	}
	return nil
}

func (r *UserRepo) ReduceBalanceByID(ctx context.Context, id uint, amount int) error {
	err := r.db.Extract(ctx).Model(&User{}).Where("id = ?", id).Update("balance", gorm.Expr("balance - ?", amount)).Error
	if err != nil {
		return reedam.Unexpected(err).WithLog(id, amount, database.IsTx(ctx))
	}
	return nil
}

func (r *UserRepo) FindBalanceByID(ctx context.Context, id uint) (int, error) {
	var u User
	err := r.db.Extract(ctx).Select("balance").First(&u, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, reedam.UserNotFound(err)
		}
		return 0, reedam.Unexpected(err).WithLog(id, database.IsTx(ctx))
	}
	return u.Balance, nil
}

func (r *UserRepo) UpdateBalanceByID(ctx context.Context, id uint, balance int) error {
	err := r.db.Extract(ctx).Model(&User{}).Where("id = ?", id).Update("balance", balance).Error
	if err != nil {
		return reedam.Unexpected(err).WithLog(id, balance, database.IsTx(ctx))
	}
	return nil
}

func (r *UserRepo) UpdateLevelXPByID(ctx context.Context, user domain.User) error {
	err := r.db.Extract(ctx).Model(&User{}).Where("id = ?", user.ID).Updates(map[string]any{
		"level": user.Level,
		"xp":    user.XP,
	}).Error
	if err != nil {
		return reedam.Unexpected(err).WithLog(user, database.IsTx(ctx))
	}
	return nil
}
