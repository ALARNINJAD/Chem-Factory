package service

import (
	"chem-factory/internal/dto/mixer"
	"chem-factory/internal/model"
)

func (s *serviceManager) AddToMixer(m mixer.MixerAddRequest) (int, error) {

	username, err := s.auth.VerifyJWT(m.Token)
	if err != nil {
		return 0, err
	}

	mm := model.Mixer{
		Username:             username,
		FirstIngredientName:  m.FirstIngredientName,
		SecondIngredientName: m.SecondIngredientName,
		Number:               m.Number,
	}

	var id int
	if id, err = s.repository.AddToMixer(mm); err != nil {
		return 0, err
	}

	mi := model.Inventory{
		Username:     username,
		MaterialName: m.FirstIngredientName,
		Number:       m.Number,
	}

	if err = s.repository.RemoveFromInventory(mi); err != nil {
		return 0, err
	}

	mi.MaterialName = m.SecondIngredientName
	if err = s.repository.RemoveFromInventory(mi); err != nil {
		return 0, err
	}

	return id, nil
}

// func (s *serviceManager) CkeckMix(id int) error {

// 	mix, err := s.repository.ExportMixRowByID(id)
// 	if err != nil {
// 		return err
// 	}

// 	mat, err := s.repository.ExportMaterialByIngrID(mix.FirstIngredientID, mix.SecondIngredientID)
// 	if err != nil {

// 	}

// }
