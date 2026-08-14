package sqlite

import (
	database "chem-factory/internal/database/sqlite"
	"chem-factory/internal/domain"
	"chem-factory/pkg/reedam"
	"context"
	"errors"
	"time"
	"gorm.io/gorm"
)

type Mix struct {
	ID                 uint      `gorm:"primaryKey;autoIncrement;not null;unique"`
	UserID             uint      `gorm:"not null;index"`
	FirstIngredientID  *uint     `gorm:"check:first_ingredient_id <= second_ingredient_id;index"`
	SecondIngredientID *uint     `gorm:"check:first_ingredient_id <= second_ingredient_id;index"`
	Amount             int       `gorm:"not null;check:10 >= amount AND amount >= 1"`
	DateTime           time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

func (Mix) TableName() string {
	return "mixes"
}

func toDomainMix(m Mix) domain.Mix {
	var fid, sid uint
	if m.FirstIngredientID != nil {
		fid = *m.FirstIngredientID
	}
	if m.SecondIngredientID != nil {
		sid = *m.SecondIngredientID
	}
	return domain.Mix{
		ID:                 m.ID,
		UserID:             m.UserID,
		FirstIngredientID:  fid,
		SecondIngredientID: sid,
		Amount:             m.Amount,
		DateTime:           m.DateTime,
	}
}

func fromDomainMix(m domain.Mix) Mix {
	var fid, sid *uint
	if m.FirstIngredientID != 0 {
		fid = &m.FirstIngredientID
	}
	if m.SecondIngredientID != 0 {
		sid = &m.SecondIngredientID
	}
	return Mix{
		ID:                 m.ID,
		UserID:             m.UserID,
		FirstIngredientID:  fid,
		SecondIngredientID: sid,
		Amount:             m.Amount,
		DateTime:           m.DateTime,
	}
}

type MixerRepo struct{ db *database.Database }

func NewMixerRepo(db *database.Database) *MixerRepo { return &MixerRepo{db: db} }

func (r *MixerRepo) FindByID(ctx context.Context, id uint) (domain.Mix, error) {
	var mix Mix
	err := r.db.Extract(ctx).First(&mix, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Mix{}, reedam.MixNotFound(err)
		}
		return domain.Mix{}, reedam.Unexpected(err).WithLog(id, database.IsTx(ctx))
	}
	return toDomainMix(mix), nil
}

func (r *MixerRepo) Add(ctx context.Context, mix domain.Mix) error {
	m := fromDomainMix(mix)
	err := r.db.Extract(ctx).Create(&m).Error
	if err != nil {
		return reedam.Unexpected(err).WithLog(mix, database.IsTx(ctx))
	}
	return nil
}

func (r *MixerRepo) FindIDByUserIDIngrID(ctx context.Context, userID, firstID, secID uint) (uint, error) {
	var mix Mix
	err := r.db.Extract(ctx).Select("id").Where("user_id = ? AND first_ingredient_id = ? AND second_ingredient_id = ?", userID, firstID, secID).First(&mix).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, reedam.MixNotFound(err)
		}
		return 0, reedam.Unexpected(err).WithLog(userID, firstID, secID, database.IsTx(ctx))
	}
	return mix.ID, nil
}

func (r *MixerRepo) GetByUserID(ctx context.Context, userID uint) ([]domain.Mix, error) {
	var mixes []Mix
	err := r.db.Extract(ctx).Where("user_id = ?", userID).Find(&mixes).Error
	if err != nil {
		return nil, reedam.Unexpected(err).WithLog(userID, database.IsTx(ctx))
	}
	list := make([]domain.Mix, len(mixes))
	for i, m := range mixes {
		list[i] = toDomainMix(m)
	}
	return list, nil
}

func (r *MixerRepo) FindByUserIDIngrID(ctx context.Context, userID, firstID, secID uint) (domain.Mix, error) {
	var mix Mix
	err := r.db.Extract(ctx).Where("user_id = ? AND first_ingredient_id = ? AND second_ingredient_id = ?", userID, firstID, secID).First(&mix).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Mix{}, reedam.MixNotFound(err)
		}
		return domain.Mix{}, reedam.Unexpected(err).WithLog(userID, firstID, secID, database.IsTx(ctx))
	}
	return toDomainMix(mix), nil
}

func (r *MixerRepo) Get(ctx context.Context) ([]domain.Mix, error) {
	var mixes []Mix
	err := r.db.Extract(ctx).Find(&mixes).Error
	if err != nil {
		return nil, reedam.Unexpected(err).WithLog(database.IsTx(ctx))
	}
	list := make([]domain.Mix, len(mixes))
	for i, m := range mixes {
		list[i] = toDomainMix(m)
	}
	return list, nil
}

func (r *MixerRepo) DeleteByID(ctx context.Context, id uint) error {
	err := r.db.Extract(ctx).Delete(&Mix{}, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return reedam.MixNotFound(err)
		}
		return reedam.Unexpected(err).WithLog(id, database.IsTx(ctx))
	}
	return nil
}
