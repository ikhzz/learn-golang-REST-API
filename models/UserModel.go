package models

type User struct {
	Id interface{} 	`bson:"_id,omitempty"` 
	Name string 		`validate:"required" bson:"name"`
	Email string 		`validate:"required,email" bson:"email"`
	Password string `validate:"required,gte=8" bson:"password"`
	Role string
}