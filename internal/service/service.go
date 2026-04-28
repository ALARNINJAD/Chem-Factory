package service

import (
	"chem-factory/internal/auth"
	i "chem-factory/internal/dto/inventory"
	"chem-factory/internal/dto/mixer"
	"chem-factory/internal/dto/shop"
	u "chem-factory/internal/dto/user"
	"chem-factory/internal/notification"
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
	CheckMix(request mixer.CheckMixRequest) (mixer.CheckMixResponse, error)
	PickMix(request mixer.PickMixRequest) error
	PickNewMix(request mixer.PickNewMixRequest) error
}

type serviceManager struct {
	auth       auth.AuthManager
	repository repository.RepositoryManager
	notification *notification.NotificationManager
}

func Init(a auth.AuthManager, r repository.RepositoryManager, m *notification.NotificationManager) *serviceManager {
	s := serviceManager{auth: a, repository: r, notification: m}
	s.createAdminUser()
	s.createAdminMaterials()
	return &s
}
