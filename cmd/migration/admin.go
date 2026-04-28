package main

import (
	"chem-factory/internal/auth"
	"chem-factory/internal/dto/material"
	"chem-factory/internal/repository"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func createAdminUser(r repository.RepositoryManager, a auth.AuthManager) {

	_, err := r.FindIDbyUsername(os.Getenv("ADMIN_NAME"))
	if err == nil {
		return
	}

	hashedPassword, err := a.HashPassword(os.Getenv("ADMIN_PASSWORD"))
	if err != nil {
		panic(fmt.Errorf("Migration, create admin: %w ", err))
	}

	if err := r.SaveUser(os.Getenv("ADMIN_NAME"), hashedPassword); err != nil {
		panic(fmt.Errorf("Migration, create admin: %w ", err))
	}
}

func createAdminMaterials(r repository.RepositoryManager) {

	adminID, err := r.FindIDbyUsername(os.Getenv("ADMIN_NAME"))
	if err != nil {
		panic(fmt.Errorf("Migration, create admin materials: %w ", err))
	}

	data, err := os.ReadFile(filepath.Join(".", "configs", "materials.json"))
	if err != nil {
		panic(fmt.Errorf("Migration, create admin materials, read file: %w ", err))
	}

	materials := r.EmptyMaterialSlice()

	if err = json.Unmarshal(data, &materials); err != nil {
		panic(fmt.Errorf("Migration, create admin materials, unmarshal: %w ", err))
	}

	for _, m := range materials {

		m.Username = os.Getenv("ADMIN_NAME")
		m.UserID = adminID

		if m.Name == m.FirstIngredientName && m.Name == m.SecondIngredientName {
			r.SaveBaseMaterial(material.BaseMaterial{
				Name:      m.Name,
				UserID:    m.UserID,
				Username:  m.Username,
				SellPrice: m.SellPrice,
				BuyPrice:  m.BuyPrice,
			})
			continue
		}

		m.FirstIngredientID, err = r.FindIDbyMaterialName(m.FirstIngredientName)
		if err != nil {
			panic(fmt.Errorf("Migration, create admin materials: %w ", err))
		}

		m.SecondIngredientID, err = r.FindIDbyMaterialName(m.SecondIngredientName)
		if err != nil {
			panic(fmt.Errorf("Migration, create admin materials: %w ", err))
		}

		if m.FirstIngredientID > m.SecondIngredientID {
			m.FirstIngredientID, m.SecondIngredientID = m.SecondIngredientID, m.FirstIngredientID
			m.FirstIngredientName, m.SecondIngredientName = m.SecondIngredientName, m.FirstIngredientName
		}

		if err = r.SaveMaterial(m); err != nil {
			panic(fmt.Errorf("Migration, create admin materials: %w ", err))
		}
	}
}
