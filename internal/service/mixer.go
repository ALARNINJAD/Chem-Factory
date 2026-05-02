package service

import (
	"chem-factory/internal/dto/mixer"
	"errors"
	"fmt"
	"log"
	"time"
)

func (s *serviceManager) AddToMixer(request mixer.MixerAddRequest) (int, error) {

	username, err := s.auth.VerifyJWT(request.Token)
	if err != nil {
		return 0, fmt.Errorf("Service mixer, add to mixer: %w ", err)
	}

	if request.FirstIngredientName == request.SecondIngredientName {
		return 0, errors.New("Service mixer, add to mixer: ")
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

	if firstID > secID {
		firstID, secID = secID, firstID
		request.FirstIngredientName, request.SecondIngredientName = request.SecondIngredientName, request.FirstIngredientName
	}

	mxr := s.repository.EmptyMixerStruct()

	mxr.UserID = userID
	mxr.FirstIngredientID = firstID
	mxr.SecondIngredientID = secID
	mxr.Username = username
	mxr.FirstIngredientName = request.FirstIngredientName
	mxr.SecondIngredientName = request.SecondIngredientName
	mxr.Number = request.Number

	log.Println(*mxr)

	id, _ = s.repository.FindMixIDByUserIDIngrID(mxr.UserID, mxr.FirstIngredientID, mxr.SecondIngredientID)
	if id != 0 {
		return id, errors.New("Service mixer, add to mixer: already exists.")
	}

	tx, err := s.repository.Transaction()
	if err != nil {
		return 0, fmt.Errorf("Service mixer, add to mixer: %w ", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err = s.repository.AddToMixer(tx, *mxr); err != nil {
		return 0, fmt.Errorf("Service mixer, add to mixer: %w ", err)
	}

	inv := s.repository.EmptyInventoryStruct()

	if invID, err = s.repository.FindInvenIDByUserIDmatID(userID, firstID); err != nil {
		return 0, fmt.Errorf("Service mixer, add to mixer: %w ", err)
	}

	if inv, err = s.repository.FindUserInvenByID(invID); err != nil {
		return 0, fmt.Errorf("Service mixer, add to mixer: %w ", err)
	}

	if inv.Number > request.Number {
		if err = s.repository.ReduceInventoryByID(tx, invID, request.Number); err != nil {
			return 0, fmt.Errorf("Service mixer, add to mixer: %w ", err)
		}
	} else if inv.Number == request.Number {
		if err = s.repository.DeleteInventoryByID(tx, invID); err != nil {
			return 0, fmt.Errorf("Service mixer, add to mixer: %w ", err)
		}
	} else {
		return 0, errors.New("Service mixer, add to mixer: not enough material.")
	}

	if invID, err = s.repository.FindInvenIDByUserIDmatID(userID, secID); err != nil {
		return 0, fmt.Errorf("Service mixer, add to mixer: %w ", err)
	}

	if inv, err = s.repository.FindUserInvenByID(invID); err != nil {
		return 0, fmt.Errorf("Service mixer, add to mixer: %w ", err)
	}

	if inv.Number > request.Number {
		if err = s.repository.ReduceInventoryByID(tx, invID, request.Number); err != nil {
			return 0, fmt.Errorf("Service mixer, add to mixer: %w ", err)
		}
	} else if inv.Number == request.Number {
		if err = s.repository.DeleteInventoryByID(tx, invID); err != nil {
			return 0, fmt.Errorf("Service mixer, add to mixer: %w ", err)
		}
	} else {
		return 0, errors.New("Service mixer, add to mixer: not enough material.")
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("Service mixer, add to mixer, commit: %w ", err)
	}

	id, err = s.repository.FindMixIDByUserIDIngrID(mxr.UserID, mxr.FirstIngredientID, mxr.SecondIngredientID)
	if err != nil {
		return 0, fmt.Errorf("Service mixer, add to mixer: %w ", err)
	}

	return id, nil
}

func (s *serviceManager) CheckMix(request mixer.CheckMixRequest) (mixer.CheckMixResponse, error) {

	username, err := s.auth.VerifyJWT(request.Token)
	if err != nil {
		return mixer.CheckMixResponse{}, fmt.Errorf("Service mixer, check mix: %w ", err)
	}

	mix, err := s.repository.FindMixRowByID(request.ID)
	if err != nil {
		return mixer.CheckMixResponse{}, fmt.Errorf("Service mixer, check mix: %w ", err)
	}

	if mix.Username != username {
		return mixer.CheckMixResponse{}, errors.New("Service mixer, check mix: not accessable.")
	}

	mat, err := s.repository.FindMaterialByIngrID(mix.FirstIngredientID, mix.SecondIngredientID)
	var timeLimit int
	var newStatus bool
	if err != nil {
		timeLimit = 30
		newStatus = true
	} else {
		timeLimit = mat.MixTime
		newStatus = false
	}

	if t := timeLimit + int(time.Until(mix.DateTime).Seconds()); t < 0 {
		return mixer.CheckMixResponse{Time: 0, NewStatus: newStatus}, nil
	} else {
		return mixer.CheckMixResponse{Time: t, NewStatus: newStatus}, nil
	}
}

func (s *serviceManager) PickMix(request mixer.PickMixRequest) error {

	r, err := s.CheckMix(mixer.CheckMixRequest{Token: request.Token, ID: request.ID})
	if err != nil {
		return fmt.Errorf("Service mixer, pick mix: %w ", err)
	}

	if r.NewStatus || r.Time > 0 {
		return errors.New("Service mixer, pick mix: not accessable.")
	}

	username, err := s.auth.VerifyJWT(request.Token)
	if err != nil {
		return fmt.Errorf("Service mixer, pick mix: %w ", err)
	}

	mxr, err := s.repository.FindMixRowByID(request.ID)
	if err != nil {
		return fmt.Errorf("Service mixer, pick mix: %w ", err)
	}

	if mxr.Username != username {
		return errors.New("Service mixer, pick mix: not accessable.")
	}

	userID, err := s.repository.FindIDbyUsername(username)
	if err != nil {
		return fmt.Errorf("Service mixer, pick mix: %w ", err)
	}

	mat, err := s.repository.FindMaterialByIngrID(mxr.FirstIngredientID, mxr.SecondIngredientID)
	if err != nil {
		return fmt.Errorf("Service mixer, pick mix: %w ", err)
	}

	tx, err := s.repository.Transaction()
	if err != nil {
		return fmt.Errorf("Service mixer, pick mix: %w ", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if invID, err := s.repository.FindInvenIDByUserIDmatID(userID, mat.ID); err != nil {
		inv := s.repository.EmptyInventoryStruct()
		inv.UserID = userID
		inv.Username = username
		inv.MaterialID = mat.ID
		inv.MaterialName = mat.Name
		inv.Number = mxr.Number
		if err = s.repository.AddToInventory(tx, *inv); err != nil {
			return fmt.Errorf("Service mixer, pick mix: %w ", err)
		}
	} else {
		if err = s.repository.IncreaseInventoryByID(tx, invID, mxr.Number); err != nil {
			return fmt.Errorf("Service mixer, pick mix: %w ", err)
		}
	}

	if err = s.repository.DeleteMixByID(tx, mxr.ID); err != nil {
		return fmt.Errorf("Service mixer, pick mix: %w ", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("Service mixer, pick mix, commit: %w ", err)
	}

	return nil
}

func (s *serviceManager) PickNewMix(request mixer.PickNewMixRequest) error {

	r, err := s.CheckMix(mixer.CheckMixRequest{Token: request.Token, ID: request.ID})
	if err != nil {
		return fmt.Errorf("Service mixer, pick new mix: %w ", err)
	}

	if !r.NewStatus || r.Time > 0 {
		return errors.New("Service mixer, pick new mix: not accessable.")
	}

	username, err := s.auth.VerifyJWT(request.Token)
	if err != nil {
		return fmt.Errorf("Service mixer, pick new mix: %w ", err)
	}

	mxr, err := s.repository.FindMixRowByID(request.ID)
	if err != nil {
		return fmt.Errorf("Service mixer, pick new mix: %w ", err)
	}

	if mxr.Username != username {
		return errors.New("Service mixer, pick new mix: not accessable.")
	}

	userID, err := s.repository.FindIDbyUsername(username)
	if err != nil {
		return fmt.Errorf("Service mixer, pick new mix: %w ", err)
	}

	firstID, err := s.repository.FindIDbyMaterialName(mxr.FirstIngredientName)
	if err != nil {
		return fmt.Errorf("Service mixer, pick new mix: %w ", err)
	}

	secID, err := s.repository.FindIDbyMaterialName(mxr.SecondIngredientName)
	if err != nil {
		return fmt.Errorf("Service mixer, pick new mix: %w ", err)
	}

	mat := s.repository.EmptyMaterialStruct()
	mat.Name = request.Name
	mat.BuyPrice = request.Price
	mat.SellPrice = request.Price * 4 / 5
	mat.MixTime = request.MixTime
	mat.FirstIngredientID = firstID
	mat.SecondIngredientID = secID
	mat.FirstIngredientName = mxr.FirstIngredientName
	mat.SecondIngredientName = mxr.SecondIngredientName
	mat.Username = username
	mat.UserID = userID

	tx, err := s.repository.Transaction()
	if err != nil {
		return fmt.Errorf("Service mixer, pick new mix: %w ", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err = s.repository.SaveMaterial(*mat); err != nil {
		return fmt.Errorf("Service mixer, pick new mix: %w ", err)
	}

	mat.ID, err = s.repository.FindMaterialIDByIngrID(mat.FirstIngredientID, mat.SecondIngredientID)
	if err != nil {
		return fmt.Errorf("Service mixer, pick new mix: %w ", err)
	}

	if invID, err := s.repository.FindInvenIDByUserIDmatID(userID, mat.ID); err != nil {
		inv := s.repository.EmptyInventoryStruct()
		inv.UserID = userID
		inv.Username = username
		inv.MaterialID = mat.ID
		inv.MaterialName = mat.Name
		inv.Number = mxr.Number
		if err = s.repository.AddToInventory(tx, *inv); err != nil {
			return fmt.Errorf("Service mixer, pick new mix: %w ", err)
		}
	} else {
		if err = s.repository.IncreaseInventoryByID(tx, invID, mxr.Number); err != nil {
			return fmt.Errorf("Service mixer, pick new mix: %w ", err)
		}
	}

	if err = s.repository.DeleteMixByID(tx, mxr.ID); err != nil {
		return fmt.Errorf("Service mixer, pick new mix: %w ", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("Service mixer, pick new mix, commit: %w ", err)
	}

	return nil
}
