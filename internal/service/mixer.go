package service

import (
	"chem-factory/internal/dto/mixer"
	"errors"
	"fmt"
)

func (s *serviceManager) AddToMixer(request mixer.MixerAddRequest) (int, error) {

	username, err := s.auth.VerifyJWT(request.Token)
	if err != nil {
		return 0, fmt.Errorf("Service mixer, add to mixer: %w ", err)
	}

	var id, userID, firstID, secID, invID int

	if userID, err = s.repository.FindIDbyUsername(username); err != nil {
		return 0, fmt.Errorf("Service mixer, add to mixer: %w ", err)
	}

	if firstID, err = s.repository.FindIDbyMaterialName(request.FirstIngredientName); err != nil {
		return 0, fmt.Errorf("Service mixer, add to mixer: %w ", err)
	}

	if secID, err = s.repository.FindIDbyMaterialName(request.SecondIngredientName); err != nil {
		return 0, fmt.Errorf("Service mixer, add to mixer: %w ", err)
	}

	mxr := s.repository.EmptyMixerStruct()

	mxr.UserID = userID
	mxr.FirstIngredientID = firstID
	mxr.SecondIngredientID = secID
	mxr.Username = username
	mxr.FirstIngredientName = request.FirstIngredientName
	mxr.SecondIngredientName = request.SecondIngredientName
	mxr.Number = request.Number

	id, err = s.repository.FindMixIDByUserIDIngrID(mxr.UserID, mxr.FirstIngredientID, mxr.SecondIngredientID)
	if err == nil {
		return id, errors.New("Service mixer, add to mixer: already exists.")
	}

	id, err = s.repository.FindMixIDByUserIDIngrID(mxr.UserID, mxr.SecondIngredientID, mxr.FirstIngredientID)
	if err == nil {
		return id, errors.New("Service mixer, add to mixer: already exists.")
	}

	if err = s.repository.AddToMixer(*mxr); err != nil {
		return 0, fmt.Errorf("Service mixer, add to mixer: %w ", err)
	}
	
	inv := s.repository.EmptyInventoryStruct()

	if invID, err = s.repository.FindInvenIDByUserIDmatID(userID,firstID); err != nil {
		return 0, fmt.Errorf("Service mixer, add to mixer: %w ", err)
	}

	if inv, err = s.repository.FindUserInvenByID(invID); err != nil {
		return 0, fmt.Errorf("Service mixer, add to mixer: %w ", err)
	}

	if inv.Number > request.Number {
		if err = s.repository.ReduceInventoryByID(invID,request.Number); err != nil {
			return 0, fmt.Errorf("Service mixer, add to mixer: %w ", err)
		}
	} else if inv.Number == request.Number {
		if err = s.repository.DeleteInventoryByID(invID); err != nil {
			return 0, fmt.Errorf("Service mixer, add to mixer: %w ", err)
		}
	} else {
		return 0, errors.New("Service mixer, add to mixer: not enough material.")
	}

	if invID, err = s.repository.FindInvenIDByUserIDmatID(userID,secID); err != nil {
		return 0, fmt.Errorf("Service mixer, add to mixer: %w ", err)
	}

	if inv, err = s.repository.FindUserInvenByID(invID); err != nil {
		return 0, fmt.Errorf("Service mixer, add to mixer: %w ", err)
	}

	if inv.Number > request.Number {
		if err = s.repository.ReduceInventoryByID(invID,request.Number); err != nil {
			return 0, fmt.Errorf("Service mixer, add to mixer: %w ", err)
		}
	} else if inv.Number == request.Number {
		if err = s.repository.DeleteInventoryByID(invID); err != nil {
			return 0, fmt.Errorf("Service mixer, add to mixer: %w ", err)
		}
	} else {
		return 0, errors.New("Service mixer, add to mixer: not enough material.")
	}

	id, err = s.repository.FindMixIDByUserIDIngrID(mxr.UserID, mxr.FirstIngredientID, mxr.SecondIngredientID)
	if err != nil {
		return 0, fmt.Errorf("Service mixer, add to mixer: %w ", err)
	}

	return id, nil
}

// func (s *serviceManager) CkeckMix(request mixer.MixerCheckMixRequest) error {

// 	username, err := s.auth.VerifyJWT(request.Token)
// 	if err != nil {
// 		return fmt.Errorf("Service mixer, add to mixer: %w ", err)
// 	}

// 	mix, err := s.repository.ExportMixRowByID(mr.ID)
// 	if err != nil {
// 		return err
// 	}

// 	if mix.Username != username {
// 		return errors.New("Not accessable.")
// 	}

// 	mat, err := s.repository.ExportMaterialByIngrID(mix.FirstIngredientID, mix.SecondIngredientID)
// 	if err != nil {

// 		var fTime, sTime int
// 		fTime, err = s.repository.ExportMatMixTimeByID(mix.FirstIngredientID)
// 		if err != nil {
// 			return err
// 		}
// 		sTime, err = s.repository.ExportMatMixTimeByID(mix.SecondIngredientID)
// 		if err != nil {
// 			return err
// 		}

// 		if (fTime+sTime)*3/4 < mat.MixTime {
// 			return errors.New("Not accessable.")
// 		} else {
// 			return nil
// 		}
// 	}

// 	if int(time.Until(mix.DateTime).Seconds()) < mat.MixTime {
// 		return errors.New("Not accessable.")
// 	} else {
// 		return nil
// 	}
// }
