package service

import (
	"chem-factory/internal/dto/shop"
	"chem-factory/internal/model"
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
