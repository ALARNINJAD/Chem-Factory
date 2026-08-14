package sqlite

import (
	database "chem-factory/internal/database/sqlite"
	"chem-factory/internal/domain"
	"chem-factory/pkg/reedam"
	"context"
	"errors"
	"gorm.io/gorm"
)

type Material struct {
	ID                 uint   `gorm:"primaryKey;autoIncrement;not null;unique"`
	UserID             *uint  `gorm:"index"`
	FirstIngredientID  *uint  `gorm:"check:first_ingredient_id <= second_ingredient_id;index"`
	SecondIngredientID *uint  `gorm:"check:first_ingredient_id <= second_ingredient_id;index"`
	Name               string `gorm:"not null;unique"`
	Price              int    `gorm:"not null;check:price >= 0"`
	MixTime            *int   `gorm:"check:mix_time >= 0"`
}

func (Material) TableName() string {
	return "materials"
}

func toDomainMaterial(m Material) domain.Material {
	var uid, fid, sid uint
	var mtime int
	if m.UserID != nil {
		uid = *m.UserID
	}
	if m.FirstIngredientID != nil {
		fid = *m.FirstIngredientID
	}
	if m.SecondIngredientID != nil {
		sid = *m.SecondIngredientID
	}
	if m.MixTime != nil {
		mtime = *m.MixTime
	}
	return domain.Material{
		ID:                 m.ID,
		UserID:             uid,
		FirstIngredientID:  fid,
		SecondIngredientID: sid,
		Name:               m.Name,
		Price:              m.Price,
		MixTime:            mtime,
	}
}

func fromDomainMaterial(m domain.Material) Material {
	var uid, fid, sid *uint
	var mtime *int
	if m.UserID != 0 {
		uid = &m.UserID
	}
	if m.FirstIngredientID != 0 {
		fid = &m.FirstIngredientID
	}
	if m.SecondIngredientID != 0 {
		sid = &m.SecondIngredientID
	}
	if m.MixTime != 0 {
		mtime = &m.MixTime
	}
	return Material{
		ID:                 m.ID,
		UserID:             uid,
		FirstIngredientID:  fid,
		SecondIngredientID: sid,
		Name:               m.Name,
		Price:              m.Price,
		MixTime:            mtime,
	}
}

type MaterialRepo struct{ db *database.Database }

func NewMaterialRepo(db *database.Database) *MaterialRepo { return &MaterialRepo{db: db} }

func (r *MaterialRepo) FindMixTimeByID(ctx context.Context, id uint) (int, error) {
	var m Material
	err := r.db.Extract(ctx).Select("mix_time").First(&m, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, reedam.MaterialNotFound(err)
		}
		return 0, reedam.Unexpected(err).WithLog(id, database.IsTx(ctx))
	}
	if m.MixTime == nil {
		return 0, nil
	}
	return *m.MixTime, nil
}

func (r *MaterialRepo) FindIDByIngrID(ctx context.Context, firstID uint, secondID uint) (uint, error) {
	var m Material
	err := r.db.Extract(ctx).Select("id").Where("first_ingredient_id = ? AND second_ingredient_id = ?", firstID, secondID).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, reedam.MaterialNotFound(err)
		}
		return 0, reedam.Unexpected(err).WithLog(firstID, secondID, database.IsTx(ctx))
	}
	return m.ID, nil
}

func (r *MaterialRepo) FindByIngrID(ctx context.Context, firstID uint, secondID uint) (domain.Material, error) {
	var m Material
	err := r.db.Extract(ctx).Where("first_ingredient_id = ? AND second_ingredient_id = ?", firstID, secondID).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Material{}, reedam.MaterialNotFound(err)
		}
		return domain.Material{}, reedam.Unexpected(err).WithLog(firstID, secondID, database.IsTx(ctx))
	}
	return toDomainMaterial(m), nil
}

func (r *MaterialRepo) FindNameByIngrID(ctx context.Context, firstID uint, secondID uint) (string, error) {
	var m Material
	err := r.db.Extract(ctx).Select("name").Where("first_ingredient_id = ? AND second_ingredient_id = ?", firstID, secondID).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", reedam.MaterialNotFound(err)
		}
		return "", reedam.Unexpected(err).WithLog(firstID, secondID, database.IsTx(ctx))
	}
	return m.Name, nil
}

func (r *MaterialRepo) FindByID(ctx context.Context, id uint) (domain.Material, error) {
	var m Material
	err := r.db.Extract(ctx).First(&m, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Material{}, reedam.MaterialNotFound(err)
		}
		return domain.Material{}, reedam.Unexpected(err).WithLog(id, database.IsTx(ctx))
	}
	return toDomainMaterial(m), nil
}

func (r *MaterialRepo) FindIDByName(ctx context.Context, name string) (uint, error) {
	var m Material
	err := r.db.Extract(ctx).Select("id").Where("name = ?", name).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, reedam.MaterialNotFound(err)
		}
		return 0, reedam.Unexpected(err).WithLog(name, database.IsTx(ctx))
	}
	return m.ID, nil
}

func (r *MaterialRepo) FindNameByID(ctx context.Context, id uint) (string, error) {
	var m Material
	err := r.db.Extract(ctx).Select("name").First(&m, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", reedam.MaterialNotFound(err)
		}
		return "", reedam.Unexpected(err).WithLog(id, database.IsTx(ctx))
	}
	return m.Name, nil
}

func (r *MaterialRepo) Add(ctx context.Context, material domain.Material) error {
	m := fromDomainMaterial(material)
	err := r.db.Extract(ctx).Create(&m).Error
	if err != nil {
		return reedam.Unexpected(err).WithLog(material, database.IsTx(ctx))
	}
	return nil
}

func (r *MaterialRepo) FindByUserID(ctx context.Context, id uint) ([]domain.Material, error) {
	var ms []Material
	err := r.db.Extract(ctx).Where("user_id = ?", id).Find(&ms).Error
	if err != nil {
		return nil, reedam.Unexpected(err).WithLog(id, database.IsTx(ctx))
	}
	list := make([]domain.Material, len(ms))
	for i, m := range ms {
		list[i] = toDomainMaterial(m)
	}
	return list, nil
}

func (r *MaterialRepo) Get(ctx context.Context) ([]domain.Material, error) {
	var ms []Material
	err := r.db.Extract(ctx).Find(&ms).Error
	if err != nil {
		return nil, reedam.Unexpected(err).WithLog(database.IsTx(ctx))
	}
	list := make([]domain.Material, len(ms))
	for i, m := range ms {
		list[i] = toDomainMaterial(m)
	}
	return list, nil
}

func (r *MaterialRepo) All(ctx context.Context) ([]domain.Material, error) {
	return r.Get(ctx)
}

func (r *MaterialRepo) FindPriceByID(ctx context.Context, id uint) (int, error) {
	var m Material
	err := r.db.Extract(ctx).Select("price").First(&m, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, reedam.MaterialNotFound(err)
		}
		return 0, reedam.Unexpected(err).WithLog(id, database.IsTx(ctx))
	}
	return m.Price, nil
}
