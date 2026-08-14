package main

import (
	"chem-factory/internal/database/sqlite"
	"chem-factory/internal/domain"
	marketsqlite "chem-factory/internal/modules/market/adapter/sqlite"
	materialsqlite "chem-factory/internal/modules/material/adapter/sqlite"
	"chem-factory/pkg/lang"
	"chem-factory/pkg/reedam"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {
	db := sqlite.New()
	ctx := context.Background()

	createMaterials(ctx, db)
	createShopMaterials(ctx, db)
	log.Println("Database seeding completed successfully.")
}

func createMaterials(ctx context.Context, db *sqlite.Database) {
	data, err := os.ReadFile(filepath.Join(".", "pkg", "material", "default_materials.json"))
	if err != nil {
		panic(fmt.Errorf("Seed, create admin materials, read file: %w ", err))
	}

	var materials []struct {
		Name                 string `json:"name"`
		Price                int    `json:"price"`
		MixTime              int    `json:"mix_time"`
		FirstIngredientName  string `json:"first_ingredient_name"`
		SecondIngredientName string `json:"second_ingredient_name"`
	}
	if err := json.Unmarshal(data, &materials); err != nil {
		panic(fmt.Errorf("Seed, create admin materials, parse file: %w ", err))
	}

	repo := materialsqlite.NewMaterialRepo(db)

	for _, m := range materials {
		if _, err := repo.FindIDByName(ctx, m.Name); err == nil {
			continue
		} else if r := reedam.As(err); r == nil || r.ErrName != lang.ErrorMaterialNotFound {
			panic(fmt.Errorf("Seed, check material existence: %w", err))
		}

		var (
			firstID  uint
			secondID uint
		)

		if m.FirstIngredientName != "" && m.SecondIngredientName != "" {
			fID, err := repo.FindIDByName(ctx, m.FirstIngredientName)
			if err != nil {
				panic(fmt.Errorf("Seed, create admin materials, find first ingredient: %w", err))
			}
			sID, err := repo.FindIDByName(ctx, m.SecondIngredientName)
			if err != nil {
				panic(fmt.Errorf("Seed, create admin materials, find second ingredient: %w", err))
			}

			if fID > sID {
				fID, sID = sID, fID
			}
			firstID = fID
			secondID = sID
		}

		if err := repo.Add(ctx, domain.Material{
			Name:               m.Name,
			Price:              m.Price,
			MixTime:            m.MixTime,
			FirstIngredientID:  firstID,
			SecondIngredientID: secondID,
		}); err != nil {
			log.Println("Could not add material:", m.Name, "error:", err)
		}
	}
}

func createShopMaterials(ctx context.Context, db *sqlite.Database) {
	materialRepo := materialsqlite.NewMaterialRepo(db)
	marketRepo := marketsqlite.NewMarketRepo(db)

	materials, err := materialRepo.All(ctx)
	if err != nil {
		panic(fmt.Errorf("Seed, read materials for shop seed: %w", err))
	}

	for _, mat := range materials {
		if !mat.IsRaw() {
			continue
		}

		exists, err := marketRepo.CheckAdminMarketExists(ctx, mat.ID)
		if err != nil {
			panic(fmt.Errorf("Seed, check shop material existence: %w", err))
		}
		if exists {
			continue
		}

		if err := marketRepo.Add(ctx, domain.Market{
			UserID:     0,
			MaterialID: mat.ID,
			Amount:     10,
			DateTime:   time.Now(),
		}); err != nil {
			log.Println("Could not add shop material:", mat.Name, "error:", err)
		}
	}
}
