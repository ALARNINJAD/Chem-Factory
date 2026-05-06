package service

import (
	"chem-factory/internal/dto/shop"
	"errors"
	"fmt"
	"log"
	"os"
)

func (service *Manager) ItemsForSell() (shop.ItemsResponse, error) {

	r := service.repository

	shopItems, err := r.Shop.Export()
	if err != nil {
		return shop.ItemsResponse{}, fmt.Errorf("Service shop, items for sell: %w ", err)
	}

	var response shop.ItemsResponse

	for _, si := range shopItems {
		response.Items = append(response.Items, shop.Shop{
			Username:     si.Username,
			MaterialName: si.MaterialName,
			Number:       si.Number,
			Price:        si.Price,
		})
	}

	materials, err := r.Material.FindByUsername(os.Getenv("ADMIN_NAME"))
	if err != nil {
		return shop.ItemsResponse{}, fmt.Errorf("Service shop, items for sell: %w ", err)
	}

	for _, m := range materials {
		response.Items = append(response.Items, shop.Shop{
			Username:     m.Username,
			MaterialName: m.Name,
			Number:       10,
			Price:        m.Price,
		})
	}

	return response, nil
}

func (service *Manager) buyAdminMaterial(request shop.BuyRequest) error {

	a := service.auth

	username, err := a.JWT.Verify(request.Token)
	if err != nil {
		return fmt.Errorf("Service shop, buy admin material: %w ", err)
	}

	var userID, matID, invID int

	r := service.repository

	if matID, err = r.Material.FindIDByName(request.MaterialName); err != nil {
		return fmt.Errorf("Service shop, buy admin material: %w ", err)
	}

	mat, err := r.Material.FindBaseByID(matID)
	if err != nil {
		return fmt.Errorf("Service shop, buy admin material: %w ", err)
	}

	log.Println(*mat)

	if mat.Price != request.Price || mat.Username != request.SellerUsername {
		return errors.New("Service shop, buy admin material: material does not exist.")
	}

	if userID, err = r.User.FindIDbyUsername(username); err != nil {
		return fmt.Errorf("Service shop, buy admin material: %w ", err)
	}

	tx, err := r.Transaction()
	if err != nil {
		return fmt.Errorf("Service user, buy admin material: %w ", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if invID, err = r.Inventory.FindIDByUserIDmatID(userID, matID); err != nil {

		inv := r.Inventory.EmptyStruct()

		inv.UserID = userID
		inv.Username = username
		inv.MaterialID = matID
		inv.MaterialName = request.MaterialName
		inv.Number = request.Number

		if err = r.Inventory.Add(tx, *inv); err != nil {
			return fmt.Errorf("Service shop, buy admin material: %w ", err)
		}

	} else {

		if err = r.Inventory.IncreaseByID(tx, invID, request.Number); err != nil {
			return fmt.Errorf("Service shop, buy admin material: %w ", err)
		}
	}

	if err = r.User.ReduceBalance(tx, username, request.Number*request.Price); err != nil {
		return fmt.Errorf("Service shop, buy admin material: %w ", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("Service shop, buy admin material, commit: %w ", err)
	}

	return nil
}

func (service *Manager) buyUserMaterial(request shop.BuyRequest) error {

	a := service.auth

	username, err := a.JWT.Verify(request.Token)
	if err != nil {
		return fmt.Errorf("Service shop, buy material: %w ", err)
	}

	var userID, materialID, shopID, invID int

	r := service.repository

	if userID, err = r.User.FindIDbyUsername(request.SellerUsername); err != nil {
		return fmt.Errorf("Service shop, buy material: %w ", err)
	}

	if materialID, err = r.Material.FindIDByName(request.MaterialName); err != nil {
		return fmt.Errorf("Service shop, buy material: %w ", err)
	}

	shopRepo, err := r.Shop.FindByInfo(userID, materialID, request.Price)
	if err != nil {
		return fmt.Errorf("Service shop, buy material: %w ", err)
	}

	tx, err := r.Transaction()
	if err != nil {
		return fmt.Errorf("Service user, buy admin material: %w ", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if shopRepo.Number > request.Number {
		if err = r.Shop.ReduceNumberByID(tx, shopID, request.Number); err != nil {
			return fmt.Errorf("Service shop, buy material: %w ", err)
		}
	} else if shopRepo.Number == request.Number {
		if err = r.Shop.DeleteByID(tx, shopID); err != nil {
			return fmt.Errorf("Service shop, buy material: %w ", err)
		}
	} else {
		return errors.New("Service shop, buy material: not anough shop items.")
	}

	if err = r.User.IncreaseBalance(tx, request.SellerUsername, request.Number*request.Price); err != nil {
		return fmt.Errorf("Service shop, buy material: %w ", err)
	}

	if userID, err = r.User.FindIDbyUsername(username); err != nil {
		return fmt.Errorf("Service shop, buy material: %w ", err)
	}

	if invID, err = r.Inventory.FindIDByUserIDmatID(userID, materialID); err != nil {

		inv := r.Inventory.EmptyStruct()

		inv.UserID = userID
		inv.Username = username
		inv.MaterialID = materialID
		inv.MaterialName = request.MaterialName
		inv.Number = request.Number

		if err = r.Inventory.Add(tx, *inv); err != nil {
			return fmt.Errorf("Service shop, buy material: %w ", err)
		}

	} else {

		if err = r.Inventory.IncreaseByID(tx, invID, request.Number); err != nil {
			return fmt.Errorf("Service shop, buy material: %w ", err)
		}
	}

	if err = r.User.ReduceBalance(tx, username, request.Number*request.Price); err != nil {
		return fmt.Errorf("Service shop, buy material: %w ", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("Service shop, buy material, commit: %w ", err)
	}

	return nil
}

func (service *Manager) BuyMaterial(request shop.BuyRequest) error {
	if request.SellerUsername == os.Getenv("ADMIN_NAME") {
		return service.buyAdminMaterial(request)
	} else {
		return service.buyUserMaterial(request)
	}
}

func (service *Manager) SetForSell(request shop.SetForSellRequest) error {

	a := service.auth

	username, err := a.JWT.Verify(request.Token)
	if err != nil {
		return fmt.Errorf("Service shop, set for sell: %w ", err)
	}

	var userID, matID, invID int

	r := service.repository

	if userID, err = r.User.FindIDbyUsername(username); err != nil {
		return fmt.Errorf("Service shop, set for sell: %w ", err)
	}

	if matID, err = r.Material.FindIDByName(request.MaterialName); err != nil {
		return fmt.Errorf("Service shop, set for sell: %w ", err)
	}

	if invID, err = r.Inventory.FindIDByUserIDmatID(userID, matID); err != nil {
		return fmt.Errorf("Service shop, set for sell: %w ", err)
	}

	inv, err := r.Inventory.FindByID(invID)
	if err != nil {
		return fmt.Errorf("Service shop, set for sell: %w ", err)
	}

	tx, err := r.Transaction()
	if err != nil {
		return fmt.Errorf("Service shop, set for sell: %w ", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if inv.Number > request.Number {
		if err = r.Inventory.ReduceByID(tx, invID, request.Number); err != nil {
			return fmt.Errorf("Service shop, set for sell: %w ", err)
		}
	} else if inv.Number == request.Number {
		if err = r.Inventory.DeleteByID(tx, invID); err != nil {
			return fmt.Errorf("Service shop, set for sell: %w ", err)
		}
	} else {
		return errors.New("Service shop, set for sell: not enough inventory items.")
	}

	shopRepo := r.Shop.EmptyStruct()

	shopRepo.MaterialID = matID
	shopRepo.UserID = userID
	shopRepo.MaterialName = request.MaterialName
	shopRepo.Username = username
	shopRepo.Number = request.Number
	shopRepo.Price = request.Price

	if err = r.Shop.Add(tx, *shopRepo); err != nil {
		return fmt.Errorf("Service shop, set for sell: %w ", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("Service shop, set for sell, commit: %w ", err)
	}

	return nil
}
