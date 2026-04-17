package service

import (
	i "chem-factory/internal/dto/inventory"
	"chem-factory/internal/model"
)

func (s *serviceManager) AddToInventory(inv i.InventoryAddRequest) (i.InventoryAddResponse, error) {

	username, err := s.auth.VerifyJWT(inv.Token)
	if err != nil {
		return i.InventoryAddResponse{}, err
	}

	err = s.repository.AddToInventory(model.Inventory{
		Username: username, MaterialName: inv.Name, Number: inv.Number})
	if err != nil {
		return i.InventoryAddResponse{}, err
	}

	return i.InventoryAddResponse{Massage: "Done."}, nil
}

func (s *serviceManager) RemoveFromInventory(inv i.InventoryAddRequest) (i.InventoryRemoveResponse, error) {

	username, err := s.auth.VerifyJWT(inv.Token)
	if err != nil {
		return i.InventoryRemoveResponse{}, err
	}

	err = s.repository.RemoveFromInventory(model.Inventory{
		Username: username, MaterialName: inv.Name, Number: inv.Number})
	if err != nil {
		return i.InventoryRemoveResponse{}, err
	}

	return i.InventoryRemoveResponse{Massage: "Done."}, nil
}

func (s *serviceManager) ExportUserInventory(inv i.InventoryExportRequest) (i.InventoryExportResponse, error) {

	username, err := s.auth.VerifyJWT(inv.Token)
	if err != nil {
		return i.InventoryExportResponse{}, err
	}

	list, err := s.repository.ExportInventory(username)
	if err != nil {
		return i.InventoryExportResponse{}, err
	}

	var ier i.InventoryExportResponse

	for _, l := range list {
		ier.InventoryList = append(ier.InventoryList, model.Inventory{
			Username: l.Username, MaterialName: l.MaterialName, Number: l.Number})
	}

	return ier, nil
}
