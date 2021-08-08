package services

import (
	"context"
	"fmt"
	model "restTestOne/models"
	"restTestOne/response"

	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductService struct {}

var trans TransactionService
// method to get all product
func (p *ProductService) GetAll() ([]model.Product, error){
	// prep variable and query all product
	var pr []model.Product
	cursor, errFind := conn.Collection("product").Find(context.Background(), bson.D{})
	if errFind != nil {
		fmt.Println(errFind)
		return pr, errFind
	}
	// iterate all cursor and append it to variable
	for cursor.Next(context.TODO()) {
    var elem model.Product
    if err := cursor.Decode(&elem); err != nil {
			fmt.Println(err)
      return pr, err
    }
    pr = append(pr, elem)
	}

	return pr, nil
}
// method to create new product
func (p *ProductService) Create(pr model.Product) (model.Product, error) {
	// validate struct
	errValid := validator.New().Struct(pr)
	if errValid != nil {
		fmt.Println(errValid)
		return pr, errValid
	}
	// query insert and check of error
	id, errInsert := conn.Collection("product").InsertOne(context.Background(),pr)
	if errInsert != nil {
		fmt.Println("err insert", errInsert)
		return pr, errInsert
	}
	idSuplier, _ := primitive.ObjectIDFromHex(pr.SuplierId)
	// set and query create transaction
	err := trans.CreateTrans(
		id.InsertedID.(primitive.ObjectID), 
		idSuplier, 
		"CREATE", 
		pr.ProductStock, 
		pr.ProductPrice * pr.ProductStock )
	if err != nil {
		fmt.Println("failed to create transaction")
	}

	return pr, nil
}
// method to update a product
func (p *ProductService) Update(pr model.Product, s string) (model.Product, error){
	// validate struct
	errValid := validator.New().Struct(pr)
	if errValid != nil {
		fmt.Println(errValid)
		return pr, errValid
	}

	id, errParse := primitive.ObjectIDFromHex(s)
	if errParse != nil {
		fmt.Println(errParse)
	}
	// only update name and price, stock will update in add function
	_, errUpdate := conn.
		Collection("product").
		UpdateOne(context.Background(), 
		bson.M{"_id": id}, 
		bson.M{"$set": bson.M{ 
			"productName": pr.ProductName,
			"productPrice": pr.ProductPrice,
		}})

	if errUpdate != nil {
		fmt.Println("err update", errUpdate)
		return pr, errUpdate
	}
	
	return pr, nil
}
// method to add or subtract/buy product
func (p *ProductService) BuyOrAdd(pr model.EditProduct) (response.ResultResponse, error) {
	// prep variable and validate struct
	var res response.ResultResponse
	var pro model.Product
	errValid := validator.New().Struct(pr)
	if errValid != nil {
		fmt.Println(errValid)
		res.Status = "failed"
		return res, errValid
	}

	proId, err := primitive.ObjectIDFromHex(pr.ProductId)
	if err != nil {
		fmt.Println(errValid)
		res.Status = "failed"
		return res, fmt.Errorf("id is not valid")	
	}
	// query find product
	errFind := conn.Collection("product").FindOne(context.Background(), bson.M{"_id" : proId}).Decode(&pro)
	if errFind != nil {
		fmt.Println(errValid)
		res.Status = "failed"
		return res, fmt.Errorf("product is not valid")	
	}
	// check if product available and equal or more than requested
	if pro.ProductStock <= 0 || pro.ProductStock < pr.Amount {
		res.Status = "failed"
		return res, fmt.Errorf("not enough stock")	
	}
	// set request buy or add
	errBuyorAdd := checkBuyOrAdd(pr, proId, pro.ProductStock, pro.ProductPrice)	
	if errBuyorAdd != nil {
		return res, fmt.Errorf("failed to update")
	}
	// set response
	res = response.ResultResponse{
		Status:"Success", 
		Amount: pr.Amount, 
		Total: pro.ProductPrice * pr.Amount,
		ProductName: pro.ProductName,
	}
	
	return res, nil
}
// method query add or substract product and detail transaction
func checkBuyOrAdd(pr model.EditProduct,id primitive.ObjectID, stock int, price int) (error) {
	// prep variable result and type of transaction
	var result int
	if pr.Type == "BUY" {
		result = stock - pr.Amount
	} else if pr.Type == "ADD" {
		result = stock + pr.Amount
	}
	// query operation for update item
	_, errUpdate := conn.
		Collection("product").
		UpdateOne(context.Background(), 
		bson.M{"_id": id}, 
		bson.M{"$set": bson.M{ 
			"productStock": result,
		}})

	if errUpdate != nil {
		fmt.Println(errUpdate)
		return fmt.Errorf("failed to update")	
	}

	oid, _ := primitive.ObjectIDFromHex(pr.CreatedBy)
	// query transaction of add or buy
	err := trans.CreateTrans(id, oid, pr.Type, pr.Amount, pr.Amount * price)
	if err != nil {
		fmt.Println("failed to create transaction")
	}
	return nil
}