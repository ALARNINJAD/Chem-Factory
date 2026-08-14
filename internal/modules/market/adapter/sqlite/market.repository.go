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

type Market struct {
	ID         uint      `gorm:"primaryKey;autoIncrement;not null;unique"`
	UserID     *uint     `gorm:"index"`
	MaterialID uint      `gorm:"not null;index"`
	Amount     int       `gorm:"not null;check:amount >= 1"`
	DateTime   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

func (Market) TableName() string {
	return "market"
}

func toDomainMarket(m Market) domain.Market {
	var uid uint
	if m.UserID != nil {
		uid = *m.UserID
	}
	return domain.Market{
		ID:         m.ID,
		UserID:     uid,
		MaterialID: m.MaterialID,
		Amount:     m.Amount,
		DateTime:   m.DateTime,
	}
}

func fromDomainMarket(m domain.Market) Market {
	var uid *uint
	if m.UserID != 0 {
		uid = &m.UserID
	}
	return Market{
		ID:         m.ID,
		UserID:     uid,
		MaterialID: m.MaterialID,
		Amount:     m.Amount,
		DateTime:   m.DateTime,
	}
}

type MarketRepo struct{ db *database.Database }

func NewMarketRepo(db *database.Database) *MarketRepo { return &MarketRepo{db: db} }

func (r *MarketRepo) Export(ctx context.Context) ([]domain.Market, error) {
	var ms []Market
	err := r.db.Extract(ctx).Find(&ms).Error
	if err != nil {
		return nil, reedam.Unexpected(err).WithLog(database.IsTx(ctx))
	}
	list := make([]domain.Market, len(ms))
	for i, m := range ms {
		list[i] = toDomainMarket(m)
	}
	return list, nil
}

func (r *MarketRepo) FindIDByUserIDMatID(ctx context.Context, userID, materialID uint) (uint, error) {
	var m Market
	err := r.db.Extract(ctx).Select("id").Where("user_id = ? AND material_id = ?", userID, materialID).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, reedam.MarketNotFound(err)
		}
		return 0, reedam.Unexpected(err).WithLog(userID, materialID, database.IsTx(ctx))
	}
	return m.ID, nil
}

func (r *MarketRepo) FindByID(ctx context.Context, id uint) (domain.Market, error) {
	var m Market
	err := r.db.Extract(ctx).First(&m, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Market{}, reedam.MarketNotFound(err)
		}
		return domain.Market{}, reedam.Unexpected(err).WithLog(id, database.IsTx(ctx))
	}
	return toDomainMarket(m), nil
}

func (r *MarketRepo) ReduceAmountByID(ctx context.Context, id uint, amount int) error {
	err := r.db.Extract(ctx).Model(&Market{}).Where("id = ?", id).Update("amount", gorm.Expr("amount - ?", amount)).Error
	if err != nil {
		return reedam.Unexpected(err).WithLog(id, amount, database.IsTx(ctx))
	}
	return nil
}

func (r *MarketRepo) IncreaseAmountByID(ctx context.Context, id uint, amount int) error {
	err := r.db.Extract(ctx).Model(&Market{}).Where("id = ?", id).Update("amount", gorm.Expr("amount + ?", amount)).Error
	if err != nil {
		return reedam.Unexpected(err).WithLog(id, amount, database.IsTx(ctx))
	}
	return nil
}

func (r *MarketRepo) Add(ctx context.Context, market domain.Market) error {
	m := fromDomainMarket(market)
	err := r.db.Extract(ctx).Create(&m).Error
	if err != nil {
		return reedam.Unexpected(err).WithLog(market, database.IsTx(ctx))
	}
	return nil
}

func (r *MarketRepo) DeleteByID(ctx context.Context, id uint) error {
	err := r.db.Extract(ctx).Delete(&Market{}, id).Error
	if err != nil {
		return reedam.Unexpected(err).WithLog(id, database.IsTx(ctx))
	}
	return nil
}

func (r *MarketRepo) FindUserIDByID(ctx context.Context, id uint) (uint, error) {
	var m Market
	err := r.db.Extract(ctx).Select("user_id").First(&m, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, reedam.MarketNotFound(err)
		}
		return 0, reedam.Unexpected(err).WithLog(id, database.IsTx(ctx))
	}
	if m.UserID == nil {
		return 0, nil
	}
	return *m.UserID, nil
}

func (r *MarketRepo) FindMatIDByID(ctx context.Context, id uint) (uint, error) {
	var m Market
	err := r.db.Extract(ctx).Select("material_id").First(&m, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, reedam.MarketNotFound(err)
		}
		return 0, reedam.Unexpected(err).WithLog(id, database.IsTx(ctx))
	}
	return m.MaterialID, nil
}

func (r *MarketRepo) CheckAdminMarketExists(ctx context.Context, materialID uint) (bool, error) {
	var count int64
	err := r.db.Extract(ctx).Model(&Market{}).Where("user_id IS NULL AND material_id = ?", materialID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
