package service

import (
	i "chem-factory/internal/dto/inventory"
	"fmt"
)

func (m *Manager) ExportUserInventory(inv i.InventoryExportRequest) (i.InventoryExportResponse, error) {

	a := m.auth

	username, err := a.JWT.Verify(inv.Token)
	if err != nil {
		return i.InventoryExportResponse{}, fmt.Errorf("Service inventory, export user inventory: %w", err)
	}

	r := m.repository

	list, err := r.Inventory.FindByUsername(username)
	if err != nil {
		return i.InventoryExportResponse{}, fmt.Errorf("Service inventory, export user inventory: %w", err)
	}

	var response i.InventoryExportResponse

	for _, l := range list {
		response.InventoryList = append(response.InventoryList, i.Inventory{
			Username: l.Username, MaterialName: l.MaterialName, Number: l.Number})
	}

	return response, nil
}
