package models

type EditProduct struct {
	ProductId string	`validate:"required"`
	Amount int 				`validate:"required,gt=0"`
	CreatedBy string 	`validate:"required"`
	Type string 			`validate:"required"`
}