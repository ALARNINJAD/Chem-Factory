package model

type Inventory struct {
	Username     string `json:"username"`
	MaterialName string `json:"material_name"`
	Number       int    `json:"number"`
}
