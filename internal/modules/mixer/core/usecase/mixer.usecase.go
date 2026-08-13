package usecase

import (
	"chem-factory/internal/domain"
	"chem-factory/internal/modules/mixer/adapter/http/dto"
	"chem-factory/internal/modules/mixer/core/port"
	"chem-factory/pkg/lang"
	"chem-factory/pkg/reedam"
	"context"
	"errors"
	"log"
	"net/http"
	"time"
)

type mixerUsecase struct {
	mixerRepo     port.MixerRepository
	userRepo      port.UserRepository
	materialRepo  port.MaterialRepository
	inventoryRepo port.InventoryRepository
	transactor    port.Transactor
}

func NewMixerUsecase(
	mixerRepo port.MixerRepository,
	userRepo port.UserRepository,
	materialRepo port.MaterialRepository,
	inventoryRepo port.InventoryRepository,
	transactor port.Transactor,
) *mixerUsecase {
	log.Println("Initializing mixer usecase")
	return &mixerUsecase{
		mixerRepo:     mixerRepo,
		userRepo:      userRepo,
		materialRepo:  materialRepo,
		inventoryRepo: inventoryRepo,
		transactor:    transactor,
	}
}

func (service *mixerUsecase) Mixes(ctx context.Context, userID uint) (dto.MixesResponse, error) {

	mixes, err := service.mixerRepo.GetByUserID(ctx, userID)
	if err != nil {
		return dto.MixesResponse{}, err
	}

	var response dto.MixesResponse
	for _, mix := range mixes {
		mixResponse := dto.MixResponse{
			ID:                 mix.ID,
			UserID:             mix.UserID,
			FirstIngredientID:  mix.FirstIngredientID,
			SecondIngredientID: mix.SecondIngredientID,
			Amount:             mix.Amount,
			DateTime:           mix.DateTime,
		}

		mixResponse.Username, err = service.userRepo.FindUsernameByID(ctx, userID)
		if err != nil {
			return dto.MixesResponse{}, err
		}

		mixResponse.FirstIngredientName, err = service.materialRepo.FindNameByID(ctx, mixResponse.FirstIngredientID)
		if err != nil {
			return dto.MixesResponse{}, err
		}

		mixResponse.SecondIngredientName, err = service.materialRepo.FindNameByID(ctx, mixResponse.SecondIngredientID)
		if err != nil {
			return dto.MixesResponse{}, err
		}

		material, _ := service.materialRepo.FindByIngrID(ctx, mix.FirstIngredientID, mix.SecondIngredientID)
		mixResponse.MaterialName = material.Name
		if material.ID == 0 {
			mixResponse.IsNew = true
			mixResponse.RemainingSeconds = mix.RemainingSeconds(30) // must be in config
		} else {
			mixResponse.IsNew = false
			mixResponse.RemainingSeconds = mix.RemainingSeconds(material.MixTime)
		}

		response.Mixes = append(response.Mixes, mixResponse)
	}

	return response, nil
}

func (service *mixerUsecase) Mix(ctx context.Context, request dto.MixRequest, userID uint) (dto.MixResponse, error) {

	if request.FirstIngredientID > request.SecondIngredientID {
		request.FirstIngredientID, request.SecondIngredientID = request.SecondIngredientID, request.FirstIngredientID
	}

	err := service.transactor.WithTx(ctx, func(ctx context.Context) error {

		inventory, err := service.inventoryRepo.FindByUserIDmatID(ctx, userID, request.FirstIngredientID)
		if err != nil {
			return err
		}
		if inventory.Amount < request.Amount {
			return reedam.New().WithError(errors.New("not enough items")).WithMessage("not enough items in").WithStatus(http.StatusConflict)
		}

		if inventory.Amount == request.Amount {
			if err = service.inventoryRepo.DeleteByID(ctx, inventory.ID); err != nil {
				return err
			}
		} else {
			if err = service.inventoryRepo.ReduceByID(ctx, inventory.ID, request.Amount); err != nil {
				return err
			}
		}

		inventory, err = service.inventoryRepo.FindByUserIDmatID(ctx, userID, request.SecondIngredientID)
		if err != nil {
			return err
		}
		if inventory.Amount < request.Amount {
			return reedam.New().WithError(errors.New("not enough items")).WithMessage("not enough items in").WithStatus(http.StatusConflict)
		}

		if inventory.Amount == request.Amount {
			if err = service.inventoryRepo.DeleteByID(ctx, inventory.ID); err != nil {
				return err
			}
		} else {
			if err = service.inventoryRepo.ReduceByID(ctx, inventory.ID, request.Amount); err != nil {
				return err
			}
		}

		err = service.mixerRepo.Add(ctx, domain.Mix{
			UserID:             userID,
			FirstIngredientID:  request.FirstIngredientID,
			SecondIngredientID: request.SecondIngredientID,
			Amount:             request.Amount,
			DateTime:           time.Now(),
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return dto.MixResponse{}, err
	}

	var mix domain.Mix
	mix, err = service.mixerRepo.FindByUserIDIngrID(ctx, userID, request.FirstIngredientID, request.SecondIngredientID)
	if err != nil {
		return dto.MixResponse{}, reedam.Unexpected(err).WithLog(mix)
	}

	var response dto.MixResponse
	response, err = service.mixResponse(ctx, mix, userID)
	if err != nil {
		return dto.MixResponse{}, reedam.Unexpected(err).WithLog(response)
	}

	response.Message = lang.MessageMixSuccessfully
	return response, nil
}

func (service *mixerUsecase) Check(ctx context.Context, request dto.CheckRequest, userID uint) (dto.MixResponse, error) {

	var (
		err      error
		mix      domain.Mix
		response dto.MixResponse
	)

	mix, err = service.mixerRepo.FindByID(ctx, request.ID)
	if err != nil {
		return dto.MixResponse{}, err
	}

	response, err = service.mixResponse(ctx, mix, userID)
	if err != nil {
		return dto.MixResponse{}, err
	}

	return response, nil
}

func (service *mixerUsecase) Pick(ctx context.Context, request dto.PickRequest, userID uint) (dto.PickResponse, error) {

	var (
		err      error
		material domain.Material
		mix      domain.Mix
		response dto.PickResponse
	)
	response.IsPicked = false

	mix, err = service.mixerRepo.FindByID(ctx, request.ID)
	if err != nil {
		return response, err
	}

	material, _ = service.materialRepo.FindByIngrID(ctx, mix.FirstIngredientID, mix.SecondIngredientID)
	if material.ID == 0 {
		response.IsNew = true
		response.Message = lang.MessagePickNewMaterial
		return response, nil
	} else {
		response.IsNew = false
	}

	response.RemainingSeconds = mix.RemainingSeconds(material.MixTime)
	if response.RemainingSeconds > 0 {
		response.Message = lang.MessagePickButRemaining
		return response, nil
	}

	err = service.transactor.WithTx(ctx, func(ctx context.Context) error {

		if err = service.mixerRepo.DeleteByID(ctx, mix.ID); err != nil {
			return err
		}

		inventory, _ := service.inventoryRepo.FindByUserIDmatID(ctx, userID, material.ID)
		if inventory.ID == 0 {
			err = service.inventoryRepo.Add(ctx, domain.Inventory{
				UserID:     userID,
				MaterialID: material.ID,
				Amount:     mix.Amount,
			})
			if err != nil {
				return err
			}
		} else {
			err = service.inventoryRepo.IncreaseByID(ctx, inventory.ID, mix.Amount)
			if err != nil {
				return err
			}
		}

		user, err := service.userRepo.FindByID(ctx, userID)
		if err != nil {
			return err
		}
		user.GetXP(material.MixTime)
		err = service.userRepo.UpdateLevelXPByID(ctx, user)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return dto.PickResponse{}, err
	}

	response.IsPicked = true
	response.Message = lang.MessagePickSuccessfully
	return response, nil
}

func (service *mixerUsecase) NewMaterial(ctx context.Context, request dto.NewMaterialRequest, userID uint) (dto.MixResponse, error) {

	var (
		mix      domain.Mix
		user     domain.User
		err      error
		response dto.MixResponse
	)
	if request.Name == "" || len(request.Name) > 50 { // these are must be inside config
		return dto.MixResponse{}, reedam.New().WithError(errors.New("invalid material name")).WithMessage("material name is not valid").WithStatus(http.StatusBadRequest)
	}
	if request.Price < 0 || request.Price > 1000000 {
		return dto.MixResponse{}, reedam.New().WithError(errors.New("invalid material price")).WithMessage("material price is not valid").WithStatus(http.StatusBadRequest)
	}
	if request.MixTime <= 0 || request.MixTime > 3600 {
		return dto.MixResponse{}, reedam.New().WithError(errors.New("invalid material mix time")).WithMessage("material mix time is not valid").WithStatus(http.StatusBadRequest)
	}

	mix, err = service.mixerRepo.FindByID(ctx, request.MixID)
	if err != nil {
		return dto.MixResponse{}, err
	}

	user, err = service.userRepo.FindByID(ctx, userID)
	if err != nil {
		return dto.MixResponse{}, err
	}

	err = service.transactor.WithTx(ctx, func(ctx context.Context) error {
		err = service.materialRepo.Add(ctx, domain.Material{
			UserID:             userID,
			FirstIngredientID:  mix.FirstIngredientID,
			SecondIngredientID: mix.SecondIngredientID,
			Name:               request.Name,
			Price:              request.Price,
			MixTime:            request.MixTime,
		})
		if err != nil {
			return err
		}

		user.GetXP(request.MixTime * request.Price)
		err = service.userRepo.UpdateLevelXPByID(ctx, user)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return dto.MixResponse{}, err
	}

	response, err = service.mixResponse(ctx, mix, userID)
	if err != nil {
		return dto.MixResponse{}, err
	}

	response.Message = lang.MessageMaterialNamedSuccessfully
	return response, nil
}

func (service *mixerUsecase) mixResponse(ctx context.Context, mix domain.Mix, userID uint) (dto.MixResponse, error) {

	if mix.UserID != userID {
		return dto.MixResponse{}, reedam.Unexpected(errors.New("you don't have this mix")).WithLog(userID, mix.UserID)
	}

	var err error
	response := dto.MixResponse{
		ID:                 mix.ID,
		UserID:             mix.UserID,
		FirstIngredientID:  mix.FirstIngredientID,
		SecondIngredientID: mix.SecondIngredientID,
		Amount:             mix.Amount,
		DateTime:           mix.DateTime,
	}

	response.Username, err = service.userRepo.FindUsernameByID(ctx, userID)
	if err != nil {
		return dto.MixResponse{}, err
	}

	response.FirstIngredientName, err = service.materialRepo.FindNameByID(ctx, response.FirstIngredientID)
	if err != nil {
		return dto.MixResponse{}, err
	}

	response.SecondIngredientName, err = service.materialRepo.FindNameByID(ctx, response.SecondIngredientID)
	if err != nil {
		return dto.MixResponse{}, err
	}

	material, _ := service.materialRepo.FindByIngrID(ctx, mix.FirstIngredientID, mix.SecondIngredientID)
	response.MaterialName = material.Name
	if material.ID == 0 {
		response.IsNew = true
		response.RemainingSeconds = mix.RemainingSeconds(30) // must be in config
	} else {
		response.IsNew = false
		response.RemainingSeconds = mix.RemainingSeconds(material.MixTime)
	}

	return response, nil
}
