package service

import (
	"chem-factory/internal/dto/mixer"
	"chem-factory/internal/model"
	"errors"
	"time"
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

func (s *serviceManager) CkeckMix(mr mixer.MixerCheckMixRequest) error {

	username, err := s.auth.VerifyJWT(mr.Token)
	if err != nil {
		return err
	}

	mix, err := s.repository.ExportMixRowByID(mr.ID)
	if err != nil {
		return err
	}

	if mix.Username != username {
		return errors.New("Not accessable.")
	}

	mat, err := s.repository.ExportMaterialByIngrID(mix.FirstIngredientID, mix.SecondIngredientID)
	if err != nil {

		var fTime, sTime int
		fTime, err = s.repository.ExportMatMixTimeByID(mix.FirstIngredientID)
		if err != nil {
			return err
		}
		sTime, err = s.repository.ExportMatMixTimeByID(mix.SecondIngredientID)
		if err != nil {
			return err
		}

		if (fTime+sTime)*3/4 < mat.MixTime {
			return errors.New("Not accessable.")
		} else {
			return nil
		}
	}

	if int(time.Until(mix.DateTime).Seconds()) < mat.MixTime {
		return errors.New("Not accessable.")
	} else {
		return nil
	}
}
