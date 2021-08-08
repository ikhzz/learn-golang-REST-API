package models

type Product struct {
	Id interface{} 			`bson:"_id,omitempty"`
	ProductName string 	`validate:"required" bson:"productName"`
	ProductPrice int 		`validate:"required,min=1000" bson:"productPrice"`
	ProductStock int 		`validate:"required,gt=0" bson:"productStock"`
	SuplierId string 		`bson:"suplierId"`
}