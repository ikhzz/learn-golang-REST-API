package services

import (
	"context"
	"fmt"
	"restTestOne/models"

	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransactionService struct {}
// method create transaction
func (t *TransactionService) CreateTrans(
	proId primitive.ObjectID, 
	createId primitive.ObjectID, 
	types string, amount int, total int) error {
	// set struct from params
	var trans = models.Transaction{
		ProductId: proId, 
		CreatedBy: createId,
		TransType: types,
		Amount: amount,
		Total: total,
	}
	// validate struct
	errValid := validator.New().Struct(trans)
	if errValid != nil {
		fmt.Println(errValid)
		return fmt.Errorf("transaction is not valid")
	}
	// query operation
	_, errInsert := conn.Collection("transaction").InsertOne(context.Background(), trans)
	if errInsert != nil {
		fmt.Println(errValid)
		return fmt.Errorf("failed to create transaction")
	}

	return nil
}