package model

type Shop struct {
	Username     string `json:"username"`
	MaterialName string `json:"material_name"`
	Number       int    `json:"number"`
	Price        int    `json:"price"`
}
