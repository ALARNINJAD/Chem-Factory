package service

import (
	"chem-factory/internal/auth"
	i "chem-factory/internal/dto/inventory"
	"chem-factory/internal/dto/shop"
	u "chem-factory/internal/dto/user"
	"chem-factory/internal/repository"
)

type ServiceManager interface {
	// user
	Login(userLR u.UserLoginRequest) (string, error)
	Register(userRR u.UserRegisterRequest) error
	UserData(token string) (u.UserDataResponse, error)
	// inventory
	AddToInventory(inv i.InventoryAddRequest) (i.InventoryAddResponse, error)
	RemoveFromInventory(inv i.InventoryAddRequest) (i.InventoryRemoveResponse, error)
	ExportUserInventory(inv i.InventoryExportRequest) (i.InventoryExportResponse, error)
	// shop
	ItemsForSell() (shop.ShopItemsResponse, error)
	Buy(shp shop.ShopBuyRequest) error
	SetForSell(shp shop.ShopSetForSellRequest) error
}

type serviceManager struct {
	auth       auth.AuthManager
	repository repository.RepositoryManager
}

func Init(a auth.AuthManager, r repository.RepositoryManager) *serviceManager {
	s := serviceManager{auth: a, repository: r}
	s.createAdminUser()
	s.createAdminMaterials()
	return &s
}
