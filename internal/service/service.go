package service

import (
	"chem-factory/internal/auth"
	i "chem-factory/internal/dto/inventory"
	"chem-factory/internal/dto/mixer"
	"chem-factory/internal/dto/shop"
	u "chem-factory/internal/dto/user"
	"chem-factory/internal/repository"
)

type ServiceManager interface {
	// user
	Login(request u.UserLoginRequest) (string, error)
	Register(request u.UserRegisterRequest) error
	UserData(token string) (u.UserDataResponse, error)
	// inventory
	AddToInventory(inv i.InventoryAddRequest) error
	RemoveFromInventory(inv i.InventoryAddRequest) error
	ExportUserInventory(inv i.InventoryExportRequest) (i.InventoryExportResponse, error)
	// shop
	ItemsForSell() (shop.ShopItemsResponse, error)
	BuyMaterial(request shop.ShopBuyRequest) error
	SetForSell(request shop.ShopSetForSellRequest) error
	// mixer
	AddToMixer(m mixer.MixerAddRequest) (int, error)
	// CkeckMix(mr mixer.MixerCheckMixRequest) error
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
