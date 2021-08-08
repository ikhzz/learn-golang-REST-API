package services

import (
	"context"
	"fmt"
	"log"
	db "restTestOne/config"
	model "restTestOne/models"
	"restTestOne/utils"
	"strings"

	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {}
// db connection property
var conn = db.MongoDb()
// get all user method
func (s *UserService) GetAll() ([]model.User , error){
	// variable for all user
	var users []model.User
	// get all data in user based on filter if err print error
	cursor, err := conn.Collection("user").Find(context.Background(), bson.D{})
	if err != nil {
		println(err)
	}
	// iterate all cursor and append it to users slice
	for cursor.Next(context.TODO()) {
    elem := model.User{}
    if err := cursor.Decode(&elem); err != nil {
            log.Fatal(err)
    }
    users = append(users, elem)
	}
	// return user
	return users, nil
}
// sign in method
func (s *UserService) SignIn(login model.LoginModel) (string, error) {
	// struct validation
	errValid := validator.New().Struct(login)
	if errValid != nil {
		errmsg := validErrMsg(errValid.Error())
		return "", fmt.Errorf(errmsg)
	}
	// set variable for query result
	var user model.User
	errUser := conn.Collection("user").FindOne(context.Background(), bson.M{"email": login.Email}).Decode(&user)
	if errUser != nil {
		fmt.Println(errUser)
		return "", fmt.Errorf("user not found")
	}
	// compare incoming and stored password
	errPass := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password))
	if errPass != nil {
		fmt.Println(errPass)
		return "", fmt.Errorf("invalid password")
	}
	// format interface to string for data in token
	id := fmt.Sprintf("%v",user.Id)
	token, err := utils.TokenGenerator(id)
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}

	return token, nil
}
// signup method
func (s *UserService) SignUp(user model.User) (string, error){
	// vlidate struct
	errValid := validator.New().Struct(user)
	if errValid != nil {
		errmsg := validErrMsg(errValid.Error())
		return "", fmt.Errorf(errmsg)
	}
	// encrypt password
	encryptPass, errEncrypt := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if errEncrypt != nil {
		return "", fmt.Errorf("failed to encrypt password")
	}
	// set role and encrypted password
	user.Password = string(encryptPass)
	user.Role = "guest"
	// query operations
	result, errInsert := conn.Collection("user").InsertOne(context.Background(), user)
	if errInsert != nil {
		return "", errInsert
	}
	// get object id
	userId := result.InsertedID.(primitive.ObjectID).Hex()
	// crate token
	token, err := utils.TokenGenerator(userId)
	if err != nil {
		return "", err
	}
	return token, nil
}

// check user service
func (s *UserService) CheckUser(id string) (bool, string, string) {
	// variable for query and roles
	var user model.User
	var types string
	// parse string to object id
	ids,err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println(err)
		return false, "id is not valid", ""
	}
	// query operation
	errUser := conn.Collection("user").FindOne(context.Background(), bson.M{"_id": ids}).Decode(&user)
	if errUser != nil {
		fmt.Println(errUser)
		return false, "User is not exist", ""
	}
	// check if user role is authorized
	roles := contains(user.Role, []string{"guest", "costumer"})
	if !roles {
		return false, "user is not authorized", ""
	}
	// set type of transaction
	if user.Role == "guest" {
		// if user role isn't sufficient
		types = "ADD"
	} else if user.Role == "costumer" {
		types = "BUY"
	}
	
	return true, id, types
}
// check error message
func validErrMsg(s string) string {
	// variable helper
	validError := []string{"Name", "Password", "Email"}
	var errList []string
	var errMsg string
	// check error
	for _,v := range validError {
		if strings.Contains(s, v) {
			errList = append(errList, v)
		}
	}
	// set error message
	switch {
		case len(errList) == 1:
			errMsg = errList[0] + " is not valid"
		case len(errList) == 2:
			errMsg = errList[0] + " and " + errList[1] + " is not valid"
		case len(errList) == 3:
			errMsg = errList[0] + ", " + errList[1] + " and " +errList[2]+ " is not valid"
	}
	
	return errMsg
}
// helper function if string in array
func contains(str string, arr []string) bool{
	for _, a := range arr {
		if a == str {
			 return true
		}
	}
	return false
}