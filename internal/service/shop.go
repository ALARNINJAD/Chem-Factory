package service

import (
	"chem-factory/internal/dto/shop"
	"chem-factory/internal/model"
	"errors"
)

func (s *serviceManager) ItemsForSell() (shop.ShopItemsResponse, error) {

	shopMaterials, err := s.repository.ExportShop()
	if err != nil {
		return shop.ShopItemsResponse{}, err
	}

	var shopItems shop.ShopItemsResponse

	for _, sm := range shopMaterials {
		shopItems.Items = append(shopItems.Items, model.Shop{
			Username:     sm.Username,
			MaterialName: sm.MaterialName,
			Number:       sm.Number,
			Price:        sm.Price,
		})
	}

	baseMaterials, err := s.repository.ExportBaseMaterials()
	if err != nil {
		return shop.ShopItemsResponse{}, err
	}

	for _, bm := range baseMaterials {

		shopItems.Items = append(shopItems.Items, model.Shop{
			Username:     bm.Username,
			MaterialName: bm.Name,
			Number:       10,
			Price:        bm.BuyPrice,
		})
	}

	return shopItems, nil
}

func (s *serviceManager) Buy(shp shop.ShopBuyRequest) error {

	username, err := s.auth.VerifyJWT(shp.Token)
	if err != nil {
		return err
	}
	if username == shp.SellerUsername {
		return errors.New("Unacceptable request.")
	}

	ms := model.Shop{
		Username:     shp.SellerUsername,
		MaterialName: shp.MaterialName,
		Number:       shp.Number,
		Price:        shp.Price,
	}

	if err = s.repository.IncreaseBalance(ms.Username, ms.Number*ms.Price); err != nil {
		return err
	}

	if ms.Username != "admin" {
		if err = s.repository.RemoveFromShop(ms); err != nil {
			return err
		}
	}

	mi := model.Inventory{
		Username:     username,
		MaterialName: shp.MaterialName,
		Number:       shp.Number,
	}

	if err = s.repository.AddToInventory(mi); err != nil {
		return err
	}

	if err = s.repository.ReduceBalance(username, ms.Price*ms.Number); err != nil {
		return err
	}

	return nil
}

func (s *serviceManager) SetForSell(shp shop.ShopSetForSellRequest) error {

	username, err := s.auth.VerifyJWT(shp.Token)
	if err != nil {
		return err
	}

	ms := model.Shop{
		Username:     username,
		MaterialName: shp.MaterialName,
		Number:       shp.Number,
		Price:        shp.Price,
	}

	if err = s.repository.AddToShop(ms); err != nil {
		return err
	}

	mi := model.Inventory{
		Username:     username,
		MaterialName: shp.MaterialName,
		Number:       shp.Number,
	}

	if err = s.repository.RemoveFromInventory(mi); err != nil {
		return err
	}

	return nil
}
