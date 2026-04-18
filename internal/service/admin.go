package service

import (
	"chem-factory/internal/model"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func (s *serviceManager) createAdminUser() {

	if id, _ := s.repository.ExportIDbyUsername("admin"); id != 0 {
		return
	}

	hashedPassword, err := s.auth.HashPassword("admin")
	if err != nil {
		panic(fmt.Sprintf("Could not register admin. %s", err.Error()))
	}

	if err := s.repository.SaveNewUser("admin", hashedPassword); err != nil {
		panic(fmt.Sprintf("Could not register admin. %s", err.Error()))
	}
}

func (s *serviceManager) createAdminMaterials() {

	var materials []model.Material

	data, err := os.ReadFile(filepath.Join(".", "configs", "admin.json"))
	if err != nil {
		panic("Could not open admin file.")
	}

	// var admin model.User
	// if err = json.Unmarshal(data, &admin); err != nil {
	// 	log.Println(admin)
	// 	panic("Could not unmarshal admin file.")
	// }

	data, err = os.ReadFile(filepath.Join(".", "configs", "materials.json"))
	if err != nil {
		panic("Could not open materials file.")
	}

	if err = json.Unmarshal(data, &materials); err != nil {
		panic("Could not unmarshal materials file.")
	}

	for _, m := range materials {
		m.Username = "admin"
		if err = s.repository.SaveMaterial(m); err != nil {
			panic(fmt.Sprintf("Could not add the %s material to database. %s", m.Name, err.Error()))
		}
	}
}
