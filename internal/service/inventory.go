package service

import (
	"chem-factory/internal/dto/inventory"
	"fmt"
)

func (service *Manager) ExportUserInventory(inv inventory.ExportRequest) (inventory.ExportResponse, error) {

	a := service.auth

	username, err := a.JWT.Verify(inv.Token)
	if err != nil {
		return inventory.ExportResponse{}, fmt.Errorf("Service inventory, export user inventory: %w", err)
	}

	r := service.repository

	list, err := r.Inventory.FindByUsername(username)
	if err != nil {
		return inventory.ExportResponse{}, fmt.Errorf("Service inventory, export user inventory: %w", err)
	}

	var response inventory.ExportResponse

	for _, l := range list {
		response.InventoryList = append(response.InventoryList, inventory.Inventory{
			Username: l.Username, MaterialName: l.MaterialName, Number: l.Number})
	}

	return response, nil
}
