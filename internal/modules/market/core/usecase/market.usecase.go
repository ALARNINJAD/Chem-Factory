package usecase

import (
	"chem-factory/internal/domain"
	"chem-factory/internal/modules/market/adapter/http/dto"
	"chem-factory/internal/modules/market/core/port"
	"context"
	"errors"
	"time"
)

type marketUsecase struct {
	materialRepo  port.MaterialRepository
	userRepo      port.UserRepository
	inventoryRepo port.InventoryRepository
	marketRepo    port.MarketRepository
	transactor    port.Transactor
}

func NewMarketUsecase(
	materialRepo port.MaterialRepository,
	userRepo port.UserRepository,
	marketRepo port.MarketRepository,
	inventoryRepo port.InventoryRepository,
	transactor port.Transactor,
) *marketUsecase {
	return &marketUsecase{
		materialRepo:  materialRepo,
		userRepo:      userRepo,
		marketRepo:    marketRepo,
		inventoryRepo: inventoryRepo,
		transactor:    transactor,
	}
}

func (service *marketUsecase) Export(ctx context.Context) (dto.MarketListResponse, error) {

	marketList, err := service.marketRepo.Export(ctx)
	if err != nil {
		return dto.MarketListResponse{}, err
	}

	var response dto.MarketListResponse

	for _, market := range marketList {
		var username string

		if market.UserID != 0 {
			username, err = service.userRepo.FindUsernameByID(ctx, market.UserID)
			if err != nil {
				return dto.MarketListResponse{}, err
			}
		}

		materialName, err := service.materialRepo.FindNameByID(ctx, market.MaterialID)
		if err != nil {
			return dto.MarketListResponse{}, err
		}
		materialPrice, err := service.materialRepo.FindPriceByID(ctx, market.MaterialID)
		if err != nil {
			return dto.MarketListResponse{}, err
		}
		response.MarketList = append(response.MarketList, dto.MarketResponse{
			ID:           market.ID,
			UserID:       market.UserID,
			MaterialID:   market.MaterialID,
			Username:     username,
			MaterialName: materialName,
			Amount:       market.Amount,
			Price:        materialPrice,
			DateTime:     market.DateTime,
		})
	}

	return response, nil
}

func (service *marketUsecase) Buy(ctx context.Context, request dto.BuyRequest, userID uint) error {

	market, err := service.marketRepo.FindByID(ctx, request.MarketID)
	if err != nil {
		return err
	}

	if userID == market.UserID {
		return errors.New("you cannot buy your own item")
	}

	if market.Amount < request.Amount {
		return errors.New("not enough items in the market")
	}

	material, err := service.materialRepo.FindByID(ctx, market.MaterialID)
	if err != nil {
		return err
	}

	return service.transactor.WithTx(ctx, func(ctx context.Context) error {

		if market.UserID != 0 {
			if market.Amount == request.Amount {
				err = service.marketRepo.DeleteByID(ctx, market.ID)
				if err != nil {
					return err
				}
			} else {
				err = service.marketRepo.ReduceAmountByID(ctx, market.ID, request.Amount)
				if err != nil {
					return err
				}
			}
			err = service.userRepo.IncreaseBalanceByID(ctx, market.UserID, request.Amount*material.Price)
			if err != nil {
				return err
			}
			user, err := service.userRepo.FindByID(ctx, market.UserID)
			if err != nil {
				return err
			}
			user.GetXP(request.Amount * (material.MixTime + material.Price))
			err = service.userRepo.UpdateLevelXPByID(ctx, user)
			if err != nil {
				return err
			}
		}
		if material.UserID != 0 {
			user, err := service.userRepo.FindByID(ctx, material.UserID)
			if err != nil {
				return err
			}
			user.GetXP((request.Amount * (material.MixTime + material.Price)) / 10)
			err = service.userRepo.UpdateLevelXPByID(ctx, user)
			if err != nil {
				return err
			}
		}

		err = service.userRepo.ReduceBalanceByID(ctx, userID, request.Amount*material.Price)
		if err != nil {
			return err
		}

		inventoryID, _ := service.inventoryRepo.FindIDByUserIDmatID(ctx, userID, market.MaterialID)
		if inventoryID == 0 {
			err := service.inventoryRepo.Add(ctx, domain.Inventory{
				UserID:     userID,
				MaterialID: market.MaterialID,
				Amount:     request.Amount,
				DateTime:   time.Now(),
			})
			if err != nil {
				return err
			}
		} else {
			err := service.inventoryRepo.IncreaseByID(ctx, inventoryID, request.Amount)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (service *marketUsecase) SetForSell(ctx context.Context, request dto.SetForSellRequest, userID uint) error {

	_, err := service.materialRepo.FindNameByID(ctx, request.MaterialID)
	if err != nil {
		return err
	}

	inventory, err := service.inventoryRepo.FindByUserIDmatID(ctx, userID, request.MaterialID)
	if err != nil {
		return err
	}
	if inventory.Amount < request.Amount {
		return errors.New("not enough items")
	}

	return service.transactor.WithTx(ctx, func(ctx context.Context) error {

		if inventory.Amount == request.Amount {
			if err = service.inventoryRepo.DeleteByID(ctx, inventory.ID); err != nil {
				return err
			}
		} else {
			err = service.inventoryRepo.ReduceByID(ctx, inventory.ID, request.Amount)
			if err != nil {
				return err
			}
		}

		marketID, _ := service.marketRepo.FindIDByUserIDMatID(ctx, userID, request.MaterialID)
		if marketID == 0 {
			err = service.marketRepo.Add(ctx, domain.Market{
				UserID:     userID,
				MaterialID: request.MaterialID,
				Amount:     request.Amount,
				DateTime:   time.Now(),
			})
			if err != nil {
				return err
			}
		} else {
			err = service.marketRepo.IncreaseAmountByID(ctx, marketID, request.Amount)
			if err != nil {
				return err
			}
		}

		return nil
	})

}
