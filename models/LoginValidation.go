package models

type LoginModel struct { 
	Email string 		`validate:"required,email" bson:"email"`
	Password string `validate:"required,gte=8" bson:"password"`
}