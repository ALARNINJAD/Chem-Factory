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

type Inventory struct {
	ID         uint      `gorm:"primaryKey;autoIncrement;not null;unique"`
	UserID     uint      `gorm:"not null;index"`
	MaterialID uint      `gorm:"not null;index"`
	Amount     int       `gorm:"not null;check:amount >= 1"`
	DateTime   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

func (Inventory) TableName() string {
	return "inventory"
}

func toDomainInventory(inv Inventory) domain.Inventory {
	return domain.Inventory{
		ID:         inv.ID,
		UserID:     inv.UserID,
		MaterialID: inv.MaterialID,
		Amount:     inv.Amount,
		DateTime:   inv.DateTime,
	}
}

func fromDomainInventory(inv domain.Inventory) Inventory {
	return Inventory{
		ID:         inv.ID,
		UserID:     inv.UserID,
		MaterialID: inv.MaterialID,
		Amount:     inv.Amount,
		DateTime:   inv.DateTime,
	}
}

type InventoryRepo struct{ db *database.Database }

func NewInventoryRepo(db *database.Database) *InventoryRepo { return &InventoryRepo{db: db} }

func (r *InventoryRepo) FindByID(ctx context.Context, id uint) (domain.Inventory, error) {
	var inv Inventory
	err := r.db.Extract(ctx).First(&inv, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Inventory{}, reedam.InventoryNotFound(err)
		}
		return domain.Inventory{}, reedam.Unexpected(err).WithLog(id, database.IsTx(ctx))
	}
	return toDomainInventory(inv), nil
}

func (r *InventoryRepo) FindByUserID(ctx context.Context, userID uint) ([]domain.Inventory, error) {
	var invs []Inventory
	err := r.db.Extract(ctx).Where("user_id = ?", userID).Find(&invs).Error
	if err != nil {
		return nil, reedam.Unexpected(err).WithLog(userID, database.IsTx(ctx))
	}
	list := make([]domain.Inventory, len(invs))
	for i, inv := range invs {
		list[i] = toDomainInventory(inv)
	}
	return list, nil
}

func (r *InventoryRepo) FindByUserIDmatID(ctx context.Context, userID, materialID uint) (domain.Inventory, error) {
	var inv Inventory
	err := r.db.Extract(ctx).Where("user_id = ? AND material_id = ?", userID, materialID).First(&inv).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Inventory{}, reedam.InventoryNotFound(err)
		}
		return domain.Inventory{}, reedam.Unexpected(err).WithLog(userID, materialID, database.IsTx(ctx))
	}
	return toDomainInventory(inv), nil
}

func (r *InventoryRepo) FindIDByUserIDmatID(ctx context.Context, userID, materialID uint) (uint, error) {
	var inv Inventory
	err := r.db.Extract(ctx).Select("id").Where("user_id = ? AND material_id = ?", userID, materialID).First(&inv).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, reedam.InventoryNotFound(err)
		}
		return 0, reedam.Unexpected(err).WithLog(userID, materialID, database.IsTx(ctx))
	}
	return inv.ID, nil
}

func (r *InventoryRepo) IncreaseByID(ctx context.Context, id uint, amount int) error {
	err := r.db.Extract(ctx).Model(&Inventory{}).Where("id = ?", id).Update("amount", gorm.Expr("amount + ?", amount)).Error
	if err != nil {
		return reedam.Unexpected(err).WithLog(id, amount, database.IsTx(ctx))
	}
	return nil
}

func (r *InventoryRepo) ReduceByID(ctx context.Context, id uint, amount int) error {
	err := r.db.Extract(ctx).Model(&Inventory{}).Where("id = ?", id).Update("amount", gorm.Expr("amount - ?", amount)).Error
	if err != nil {
		return reedam.Unexpected(err).WithLog(id, amount, database.IsTx(ctx))
	}
	return nil
}

func (r *InventoryRepo) Add(ctx context.Context, inventory domain.Inventory) error {
	inv := fromDomainInventory(inventory)
	err := r.db.Extract(ctx).Create(&inv).Error
	if err != nil {
		return reedam.Unexpected(err).WithLog(inventory, database.IsTx(ctx))
	}
	return nil
}

func (r *InventoryRepo) DeleteByID(ctx context.Context, id uint) error {
	err := r.db.Extract(ctx).Delete(&Inventory{}, id).Error
	if err != nil {
		return reedam.Unexpected(err).WithLog(id, database.IsTx(ctx))
	}
	return nil
}
