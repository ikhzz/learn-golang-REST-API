package models

type Transaction struct {
	Id interface{} 				`bson:"_id,omitempty"` 
	ProductId interface{} `validate:"required" bson:"productId"`
	CreatedBy interface{} `validate:"required" bson:"transId"`
	TransType string 			`validate:"required" bson:"transType"`
	Amount int 						`validate:"required,gt=0" bson:"amount"`
	Total int 						`validate:"required,gt=0" bson:"total"`
}