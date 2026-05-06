package service

import (
	"chem-factory/internal/dto/mixer"
	"errors"
	"fmt"
	"log"
	"time"
)

func (service *Manager) AddToMixer(request mixer.AddRequest) (int, error) {

	a := service.auth

	username, err := a.JWT.Verify(request.Token)
	if err != nil {
		return 0, fmt.Errorf("Service mixer, add to mixer: %w", err)
	}

	if request.FirstIngredientName == request.SecondIngredientName {
		return 0, errors.New("Service mixer, add to mixer: same materials exist.")
	}

	var id, userID, firstID, secID, invID int

	r := service.repository

	if userID, err = r.User.FindIDbyUsername(username); err != nil {
		return 0, fmt.Errorf("Service mixer, add to mixer: %w ", err)
	}

	if firstID, err = r.Material.FindIDByName(request.FirstIngredientName); err != nil {
		return 0, fmt.Errorf("Service mixer, add to mixer: %w ", err)
	}

	if secID, err = r.Material.FindIDByName(request.SecondIngredientName); err != nil {
		return 0, fmt.Errorf("Service mixer, add to mixer: %w ", err)
	}

	if firstID > secID {
		firstID, secID = secID, firstID
		request.FirstIngredientName, request.SecondIngredientName = request.SecondIngredientName, request.FirstIngredientName
	}

	mxr := r.Mixer.EmptyStruct()

	mxr.UserID = userID
	mxr.FirstIngredientID = firstID
	mxr.SecondIngredientID = secID
	mxr.Username = username
	mxr.FirstIngredientName = request.FirstIngredientName
	mxr.SecondIngredientName = request.SecondIngredientName
	mxr.Number = request.Number

	log.Println(*mxr)

	id, _ = r.Mixer.FindIDByUserIDIngrID(mxr.UserID, mxr.FirstIngredientID, mxr.SecondIngredientID)
	if id != 0 {
		return id, errors.New("Service mixer, add to mixer: already exists.")
	}

	tx, err := r.Transaction()
	if err != nil {
		return 0, fmt.Errorf("Service mixer, add to mixer: %w ", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err = r.Mixer.Add(tx, *mxr); err != nil {
		return 0, fmt.Errorf("Service mixer, add to mixer: %w ", err)
	}

	inv := r.Inventory.EmptyStruct()

	if invID, err = r.Inventory.FindIDByUserIDmatID(userID, firstID); err != nil {
		return 0, fmt.Errorf("Service mixer, add to mixer: %w ", err)
	}

	if inv, err = r.Inventory.FindByID(invID); err != nil {
		return 0, fmt.Errorf("Service mixer, add to mixer: %w ", err)
	}

	if inv.Number > request.Number {
		if err = r.Inventory.ReduceByID(tx, invID, request.Number); err != nil {
			return 0, fmt.Errorf("Service mixer, add to mixer: %w ", err)
		}
	} else if inv.Number == request.Number {
		if err = r.Inventory.DeleteByID(tx, invID); err != nil {
			return 0, fmt.Errorf("Service mixer, add to mixer: %w ", err)
		}
	} else {
		return 0, errors.New("Service mixer, add to mixer: not enough material.")
	}

	if invID, err = r.Inventory.FindIDByUserIDmatID(userID, secID); err != nil {
		return 0, fmt.Errorf("Service mixer, add to mixer: %w ", err)
	}

	if inv, err = r.Inventory.FindByID(invID); err != nil {
		return 0, fmt.Errorf("Service mixer, add to mixer: %w ", err)
	}

	if inv.Number > request.Number {
		if err = r.Inventory.ReduceByID(tx, invID, request.Number); err != nil {
			return 0, fmt.Errorf("Service mixer, add to mixer: %w ", err)
		}
	} else if inv.Number == request.Number {
		if err = r.Inventory.DeleteByID(tx, invID); err != nil {
			return 0, fmt.Errorf("Service mixer, add to mixer: %w ", err)
		}
	} else {
		return 0, errors.New("Service mixer, add to mixer: not enough material.")
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("Service mixer, add to mixer, commit: %w ", err)
	}

	id, err = r.Mixer.FindIDByUserIDIngrID(mxr.UserID, mxr.FirstIngredientID, mxr.SecondIngredientID)
	if err != nil {
		return 0, fmt.Errorf("Service mixer, add to mixer: %w ", err)
	}

	return id, nil
}

func (service *Manager) CheckMix(request mixer.CheckRequest) (mixer.CheckResponse, error) {

	a := service.auth

	username, err := a.JWT.Verify(request.Token)
	if err != nil {
		return mixer.CheckResponse{}, fmt.Errorf("Service mixer, check mix: %w ", err)
	}

	r := service.repository

	mix, err := r.Mixer.FindByID(request.ID)
	if err != nil {
		return mixer.CheckResponse{}, fmt.Errorf("Service mixer, check mix: %w ", err)
	}

	if mix.Username != username {
		return mixer.CheckResponse{}, errors.New("Service mixer, check mix: not accessable.")
	}

	mat, err := r.Material.FindByIngrID(mix.FirstIngredientID, mix.SecondIngredientID)
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
		return mixer.CheckResponse{Time: 0, NewStatus: newStatus}, nil
	} else {
		return mixer.CheckResponse{Time: t, NewStatus: newStatus}, nil
	}
}

func (service *Manager) PickMix(request mixer.PickRequest) error {

	resp, err := service.CheckMix(mixer.CheckRequest{Token: request.Token, ID: request.ID})
	if err != nil {
		return fmt.Errorf("Service mixer, pick mix: %w ", err)
	}

	if resp.NewStatus || resp.Time > 0 {
		return errors.New("Service mixer, pick mix: not accessable.")
	}

	a := service.auth

	username, err := a.JWT.Verify(request.Token)
	if err != nil {
		return fmt.Errorf("Service mixer, pick mix: %w ", err)
	}

	r := service.repository

	mxr, err := r.Mixer.FindByID(request.ID)
	if err != nil {
		return fmt.Errorf("Service mixer, pick mix: %w ", err)
	}

	if mxr.Username != username {
		return errors.New("Service mixer, pick mix: not accessable.")
	}

	userID, err := r.User.FindIDbyUsername(username)
	if err != nil {
		return fmt.Errorf("Service mixer, pick mix: %w ", err)
	}

	mat, err := r.Material.FindByIngrID(mxr.FirstIngredientID, mxr.SecondIngredientID)
	if err != nil {
		return fmt.Errorf("Service mixer, pick mix: %w ", err)
	}

	tx, err := r.Transaction()
	if err != nil {
		return fmt.Errorf("Service mixer, pick mix: %w ", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if invID, err := r.Inventory.FindIDByUserIDmatID(userID, mat.ID); err != nil {
		inv := r.Inventory.EmptyStruct()
		inv.UserID = userID
		inv.Username = username
		inv.MaterialID = mat.ID
		inv.MaterialName = mat.Name
		inv.Number = mxr.Number
		if err = r.Inventory.Add(tx, *inv); err != nil {
			return fmt.Errorf("Service mixer, pick mix: %w ", err)
		}
	} else {
		if err = r.Inventory.IncreaseByID(tx, invID, mxr.Number); err != nil {
			return fmt.Errorf("Service mixer, pick mix: %w ", err)
		}
	}

	if err = r.Mixer.DeleteByID(tx, mxr.ID); err != nil {
		return fmt.Errorf("Service mixer, pick mix: %w ", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("Service mixer, pick mix, commit: %w ", err)
	}

	return nil
}

func (service *Manager) PickNewMix(request mixer.PickNewRequest) error {

	resp, err := service.CheckMix(mixer.CheckRequest{Token: request.Token, ID: request.ID})
	if err != nil {
		return fmt.Errorf("Service mixer, pick new mix: %w ", err)
	}

	if !resp.NewStatus || resp.Time > 0 {
		return errors.New("Service mixer, pick new mix: not accessable.")
	}

	a := service.auth

	username, err := a.JWT.Verify(request.Token)
	if err != nil {
		return fmt.Errorf("Service mixer, pick new mix: %w ", err)
	}

	r := service.repository

	mxr, err := r.Mixer.FindByID(request.ID)
	if err != nil {
		return fmt.Errorf("Service mixer, pick new mix: %w ", err)
	}

	if mxr.Username != username {
		return errors.New("Service mixer, pick new mix: not accessable.")
	}

	userID, err := r.User.FindIDbyUsername(username)
	if err != nil {
		return fmt.Errorf("Service mixer, pick new mix: %w ", err)
	}

	firstID, err := r.Material.FindIDByName(mxr.FirstIngredientName)
	if err != nil {
		return fmt.Errorf("Service mixer, pick new mix: %w ", err)
	}

	secID, err := r.Material.FindIDByName(mxr.SecondIngredientName)
	if err != nil {
		return fmt.Errorf("Service mixer, pick new mix: %w ", err)
	}

	mat := r.Material.EmptyStruct()
	mat.Name = request.Name
	mat.Price = request.Price
	mat.Price = request.Price * 4 / 5
	mat.MixTime = request.MixTime
	mat.FirstIngredientID = firstID
	mat.SecondIngredientID = secID
	mat.FirstIngredientName = mxr.FirstIngredientName
	mat.SecondIngredientName = mxr.SecondIngredientName
	mat.Username = username
	mat.UserID = userID

	tx, err := r.Transaction()
	if err != nil {
		return fmt.Errorf("Service mixer, pick new mix: %w ", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err = r.Material.Add(*mat); err != nil {
		return fmt.Errorf("Service mixer, pick new mix: %w ", err)
	}

	mat.ID, err = r.Material.FindIDByIngrID(mat.FirstIngredientID, mat.SecondIngredientID)
	if err != nil {
		return fmt.Errorf("Service mixer, pick new mix: %w ", err)
	}

	if invID, err := r.Inventory.FindIDByUserIDmatID(userID, mat.ID); err != nil {
		inv := r.Inventory.EmptyStruct()
		inv.UserID = userID
		inv.Username = username
		inv.MaterialID = mat.ID
		inv.MaterialName = mat.Name
		inv.Number = mxr.Number
		if err = r.Inventory.Add(tx, *inv); err != nil {
			return fmt.Errorf("Service mixer, pick new mix: %w ", err)
		}
	} else {
		if err = r.Inventory.IncreaseByID(tx, invID, mxr.Number); err != nil {
			return fmt.Errorf("Service mixer, pick new mix: %w ", err)
		}
	}

	if err = r.Mixer.DeleteByID(tx, mxr.ID); err != nil {
		return fmt.Errorf("Service mixer, pick new mix: %w ", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("Service mixer, pick new mix, commit: %w ", err)
	}

	return nil
}
