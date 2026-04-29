package service

import (
	"chem-factory/internal/dto/shop"
	"errors"
	"fmt"
	"log"
	"os"
)

func (s *serviceManager) ItemsForSell() (shop.ShopItemsResponse, error) {

	shopItems, err := s.repository.ExportShop()
	if err != nil {
		return shop.ShopItemsResponse{}, fmt.Errorf("Service shop, items for sell: %w ", err)
	}

	var response shop.ShopItemsResponse

	for _, si := range shopItems {
		response.Items = append(response.Items, shop.Shop{
			Username:     si.Username,
			MaterialName: si.MaterialName,
			Number:       si.Number,
			Price:        si.Price,
		})
	}

	materials, err := s.repository.GetMaterialsByUsername("admin")
	if err != nil {
		return shop.ShopItemsResponse{}, fmt.Errorf("Service shop, items for sell: %w ", err)
	}

	for _, m := range materials {
		response.Items = append(response.Items, shop.Shop{
			Username:     m.Username,
			MaterialName: m.Name,
			Number:       10,
			Price:        m.BuyPrice,
		})
	}

	return response, nil
}

func (s *serviceManager) buyAdminMaterial(request shop.ShopBuyRequest) error {

	username, err := s.auth.VerifyJWT(request.Token)
	if err != nil {
		return fmt.Errorf("Service shop, buy admin material: %w ", err)
	}

	var userID, matID, invID int

	if matID, err = s.repository.FindIDbyMaterialName(request.MaterialName); err != nil {
		return fmt.Errorf("Service shop, buy admin material: %w ", err)
	}

	mat, err := s.repository.FindBaseMaterialByID(matID)
	if err != nil {
		return fmt.Errorf("Service shop, buy admin material: %w ", err)
	}

	log.Println(*mat)

	if mat.BuyPrice != request.Price || mat.Username != request.SellerUsername {
		return errors.New("Service shop, buy admin material: material does not exist.")
	}

	if userID, err = s.repository.FindIDbyUsername(username); err != nil {
		return fmt.Errorf("Service shop, buy admin material: %w ", err)
	}

	tx, err := s.repository.Transaction()
	if err != nil {
		return fmt.Errorf("Service user, buy admin material: %w ", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if invID, err = s.repository.FindInvenIDByUserIDmatID(userID, matID); err != nil {

		inv := s.repository.EmptyInventoryStruct()

		inv.UserID = userID
		inv.Username = username
		inv.MaterialID = matID
		inv.MaterialName = request.MaterialName
		inv.Number = request.Number

		if err = s.repository.AddToInventory(tx, *inv); err != nil {
			return fmt.Errorf("Service shop, buy admin material: %w ", err)
		}

	} else {

		if err = s.repository.IncreaseInventoryByID(tx, invID, request.Number); err != nil {
			return fmt.Errorf("Service shop, buy admin material: %w ", err)
		}
	}

	if err = s.repository.ReduceBalance(tx, username, request.Number*request.Price); err != nil {
		return fmt.Errorf("Service shop, buy admin material: %w ", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("Service shop, buy admin material, commit: %w ", err)
	}

	return nil
}

func (s *serviceManager) buyUserMaterial(request shop.ShopBuyRequest) error {

	username, err := s.auth.VerifyJWT(request.Token)
	if err != nil {
		return fmt.Errorf("Service shop, buy material: %w ", err)
	}

	var userID, materialID, shopID, invID int

	if userID, err = s.repository.FindIDbyUsername(request.SellerUsername); err != nil {
		return fmt.Errorf("Service shop, buy material: %w ", err)
	}

	if materialID, err = s.repository.FindIDbyMaterialName(request.MaterialName); err != nil {
		return fmt.Errorf("Service shop, buy material: %w ", err)
	}

	shopRepo, err := s.repository.FindShopByInfo(userID, materialID, request.Price)
	if err != nil {
		return fmt.Errorf("Service shop, buy material: %w ", err)
	}

	tx, err := s.repository.Transaction()
	if err != nil {
		return fmt.Errorf("Service user, buy admin material: %w ", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if shopRepo.Number > request.Number {
		if err = s.repository.ReduceShopNumberByID(tx, shopID, request.Number); err != nil {
			return fmt.Errorf("Service shop, buy material: %w ", err)
		}
	} else if shopRepo.Number == request.Number {
		if err = s.repository.DeleteFromShopByID(tx, shopID); err != nil {
			return fmt.Errorf("Service shop, buy material: %w ", err)
		}
	} else {
		return errors.New("Service shop, buy material: not anough shop items.")
	}

	if err = s.repository.IncreaseBalance(tx, request.SellerUsername, request.Number*request.Price); err != nil {
		return fmt.Errorf("Service shop, buy material: %w ", err)
	}

	if userID, err = s.repository.FindIDbyUsername(username); err != nil {
		return fmt.Errorf("Service shop, buy material: %w ", err)
	}

	if invID, err = s.repository.FindInvenIDByUserIDmatID(userID, materialID); err != nil {

		inv := s.repository.EmptyInventoryStruct()

		inv.UserID = userID
		inv.Username = username
		inv.MaterialID = materialID
		inv.MaterialName = request.MaterialName
		inv.Number = request.Number

		if err = s.repository.AddToInventory(tx, *inv); err != nil {
			return fmt.Errorf("Service shop, buy material: %w ", err)
		}

	} else {

		if err = s.repository.IncreaseInventoryByID(tx, invID, request.Number); err != nil {
			return fmt.Errorf("Service shop, buy material: %w ", err)
		}
	}

	if err = s.repository.ReduceBalance(tx, username, request.Number*request.Price); err != nil {
		return fmt.Errorf("Service shop, buy material: %w ", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("Service shop, buy material, commit: %w ", err)
	}

	return nil
}

func (s *serviceManager) BuyMaterial(request shop.ShopBuyRequest) error {

	if request.SellerUsername == os.Getenv("ADMIN_NAME") {
		return s.buyAdminMaterial(request)
	} else {
		return s.buyUserMaterial(request)
	}
}

func (s *serviceManager) SetForSell(request shop.ShopSetForSellRequest) error {

	username, err := s.auth.VerifyJWT(request.Token)
	if err != nil {
		return fmt.Errorf("Service shop, set for sell: %w ", err)
	}

	var userID, matID, invID int

	if userID, err = s.repository.FindIDbyUsername(username); err != nil {
		return fmt.Errorf("Service shop, set for sell: %w ", err)
	}

	if matID, err = s.repository.FindIDbyMaterialName(request.MaterialName); err != nil {
		return fmt.Errorf("Service shop, set for sell: %w ", err)
	}

	if invID, err = s.repository.FindInvenIDByUserIDmatID(userID, matID); err != nil {
		return fmt.Errorf("Service shop, set for sell: %w ", err)
	}

	inv, err := s.repository.FindUserInvenByID(invID)
	if err != nil {
		return fmt.Errorf("Service shop, set for sell: %w ", err)
	}

	tx, err := s.repository.Transaction()
	if err != nil {
		return fmt.Errorf("Service shop, set for sell: %w ", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if inv.Number > request.Number {
		if err = s.repository.ReduceInventoryByID(tx, invID, request.Number); err != nil {
			return fmt.Errorf("Service shop, set for sell: %w ", err)
		}
	} else if inv.Number == request.Number {
		if err = s.repository.DeleteInventoryByID(tx, invID); err != nil {
			return fmt.Errorf("Service shop, set for sell: %w ", err)
		}
	} else {
		return errors.New("Service shop, set for sell: not enough inventory items.")
	}

	shopRepo := s.repository.EmptyShopStruct()

	shopRepo.MaterialID = matID
	shopRepo.UserID = userID
	shopRepo.MaterialName = request.MaterialName
	shopRepo.Username = username
	shopRepo.Number = request.Number
	shopRepo.Price = request.Price

	if err = s.repository.AddToShop(tx, *shopRepo); err != nil {
		return fmt.Errorf("Service shop, set for sell: %w ", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("Service shop, set for sell, commit: %w ", err)
	}

	return nil
}
