package service

import (
	i "chem-factory/internal/dto/inventory"
	"fmt"
)

func (s *serviceManager) AddToInventory(inv i.InventoryAddRequest) error {

	username, err := s.auth.VerifyJWT(inv.Token)
	if err != nil {
		return err
	}

	var matID, userID int
	inventory := s.repository.EmptyInventoryStruct()

	if matID, err = s.repository.FindIDbyMaterialName(inv.Name); err != nil {
		return fmt.Errorf("Service inventory, add to inventory: %w ", err)
	}

	if userID, err = s.repository.FindIDbyUsername(username); err != nil {
		return fmt.Errorf("Service inventory, add to inventory: %w ", err)
	}

	inventory.UserID = userID
	inventory.MaterialID = matID
	inventory.Username = username
	inventory.MaterialName = inv.Name
	inventory.Number = inv.Number

	if err = s.repository.AddToInventory(*inventory); err != nil {
		return fmt.Errorf("Service inventory, add to inventory: %w ", err)
	}

	return nil
}

func (s *serviceManager) RemoveFromInventory(inv i.InventoryAddRequest) error {

	username, err := s.auth.VerifyJWT(inv.Token)
	if err != nil {
		return fmt.Errorf("Service inventory, remove from inventory: %w ", err)
	}

	var invID, userID, matID int

	if userID, err = s.repository.FindIDbyUsername(username); err != nil {
		return fmt.Errorf("Service inventory, remove from inventory: %w ", err)
	}

	if matID, err = s.repository.FindIDbyMaterialName(inv.Name); err != nil {
		return fmt.Errorf("Service inventory, remove from inventory: %w ", err)
	}

	if invID, err = s.repository.FindInvenIDByUserIDmatID(userID, matID); err != nil {
		return fmt.Errorf("Service inventory, remove from inventory: %w ", err)
	}

	inventory, err := s.repository.FindUserInvenByID(invID)
	if err != nil {
		return fmt.Errorf("Service inventory, remove from inventory: %w ", err)
	}

	if inventory.Number > inv.Number {

		if err = s.repository.ReduceBalance(inventory.Username, inv.Number); err != nil {
			return fmt.Errorf("Service inventory, remove from inventory: %w ", err)
		}

	} else if inventory.Number == inv.Number {

		if err = s.repository.DeleteInventoryByID(invID); err != nil {
			return fmt.Errorf("Service inventory, remove from inventory: %w ", err)
		}

	} else {

		return fmt.Errorf("Service inventory, remove from inventory: low balance.")
	}

	return nil
}

func (s *serviceManager) ExportUserInventory(inv i.InventoryExportRequest) (i.InventoryExportResponse, error) {

	username, err := s.auth.VerifyJWT(inv.Token)
	if err != nil {
		return i.InventoryExportResponse{}, err
	}

	list, err := s.repository.GetInventoryItemsByUsername(username)
	if err != nil {
		return i.InventoryExportResponse{}, err
	}

	var response i.InventoryExportResponse

	for _, l := range list {
		response.InventoryList = append(response.InventoryList, i.Inventory{
			Username: l.Username, MaterialName: l.MaterialName, Number: l.Number})
	}

	return response, nil
}
