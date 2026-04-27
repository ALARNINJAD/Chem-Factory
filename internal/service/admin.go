package service

import (
	"chem-factory/internal/dto/material"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func (s *serviceManager) createAdminUser() {

	_, err := s.repository.FindIDbyUsername("admin")
	if err == nil {
		return
	}

	hashedPassword, err := s.auth.HashPassword("admin")
	if err != nil {
		panic(fmt.Errorf("Service admin, create admin: %w ", err))
	}

	if err := s.repository.SaveUser("admin", hashedPassword); err != nil {
		panic(fmt.Errorf("Service admin, create admin: %w ", err))
	}
}

func (s *serviceManager) createAdminMaterials() {

	adminID, err := s.repository.FindIDbyUsername("admin")
	if err != nil {
		panic(fmt.Errorf("Service admin, create admin materials: %w ", err))
	}

	data, err := os.ReadFile(filepath.Join(".", "configs", "materials.json"))
	if err != nil {
		panic(fmt.Errorf("Service admin, create admin materials, read file: %w ", err))
	}

	materials := s.repository.EmptyMaterialSlice()

	if err = json.Unmarshal(data, &materials); err != nil {
		panic(fmt.Errorf("Service admin, create admin materials, unmarshal: %w ", err))
	}

	for _, m := range materials {

		m.Username = "admin"
		m.UserID = adminID

		if m.Name == m.FirstIngredientName && m.Name == m.SecondIngredientName {
			s.repository.SaveBaseMaterial(material.BaseMaterial{
				Name:      m.Name,
				UserID:    m.UserID,
				Username:  m.Username,
				SellPrice: m.SellPrice,
				BuyPrice:  m.BuyPrice,
			})
			continue
		}

		m.FirstIngredientID, err = s.repository.FindIDbyMaterialName(m.FirstIngredientName)
		if err != nil {
			panic(fmt.Errorf("Service admin, create admin materials: %w ", err))
		}

		m.SecondIngredientID, err = s.repository.FindIDbyMaterialName(m.SecondIngredientName)
		if err != nil {
			panic(fmt.Errorf("Service admin, create admin materials: %w ", err))
		}

		_, err = s.repository.FindMaterialIDByIngrID(m.FirstIngredientID, m.SecondIngredientID)
		if err != nil {
			_, err = s.repository.FindMaterialIDByIngrID(m.SecondIngredientID, m.FirstIngredientID)
			if err != nil {
				if m.FirstIngredientID > m.SecondIngredientID {
					m.FirstIngredientID, m.SecondIngredientID = m.SecondIngredientID, m.FirstIngredientID
					m.FirstIngredientName, m.SecondIngredientName = m.SecondIngredientName, m.FirstIngredientName
				}
				if err = s.repository.SaveMaterial(m); err != nil {
					panic(fmt.Errorf("Service admin, create admin materials, %s: %w ", m.Name, err))
				}
			}
		}
	}
}
