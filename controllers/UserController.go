package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	model "restTestOne/models"
	ress "restTestOne/response"
	"restTestOne/services"
)

type UserController struct {}

var svc = services.UserService{}
// get all controller
func (controller *UserController) GetAll(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	res, err := svc.GetAll()

	if err != nil {
		json.NewEncoder(w).Encode("err")	
	}
	
	json.NewEncoder(w).Encode(res)
}
// sign in controller
func (controller *UserController) GetOne(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	// headers for cors
	// w.Header().Set("Access-Control-Allow-Origin", "*")
  // w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	var signin model.LoginModel
	// parse x-www-form-urlencoded or multipart/formdata
	r.ParseMultipartForm(0)
	for k,v := range r.Form {
		// check input form
		check, str := contains(k, []string{"email", "password"})
		if !check {
			json.NewEncoder(w).Encode(str+" is not a valid input")
			return
		}
		// set form data to struct
		if k == "email" {
			signin.Email = v[0]
		}
		if k == "password" {
			signin.Password = v[0]
		}
	}
	// parse JSON stringify or JSON object
	json.NewDecoder(r.Body).Decode(&signin)
	
	fmt.Println(signin)
	// service will validate if struct is valid
	token, err := svc.SignIn(signin)
	if err != nil {
		res := ress.ErrorResponse{Status: "Bad Request", Message: err.Error()}
		json.NewEncoder(w).Encode(res)
		return
	}

	res := ress.TokenResponse{Status: "Sign Up Success", Token: token}
	json.NewEncoder(w).Encode(res)
}
// signup controller
func (controller *UserController) Create(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	// setup user struct
	var signup model.User
	// parse x-www-form-urlencoded or multipart/formdata
	r.ParseMultipartForm(0)	
	for k,v := range r.Form {
		// check input form
		check, str := contains(k, []string{"name", "email", "password"})
		if !check {
			json.NewEncoder(w).Encode(str+" is not a valid input")
			return
		}
		// set form data to struct
		if k == "name" {
			signup.Name = v[0]
		}
		if k == "email" {
			signup.Email = v[0]
		}
		if k == "password" {
			signup.Password = v[0]
		}
	}
	// parse JSON stringify or JSON object
	json.NewDecoder(r.Body).Decode(&signup)
	// service will validate if struct is valid
	token, err := svc.SignUp(signup)
	if err != nil {
		res := ress.ErrorResponse{Status: "Bad Request", Message: err.Error()}
		json.NewEncoder(w).Encode(res)
		return
	}

	res := ress.TokenResponse{Status: "Sign Up Success", Token: token}
	json.NewEncoder(w).Encode(res)
}

func contains(str string, arr []string) (bool, string){
	for _, a := range arr {
		if a == str {
			 return true, str
		}
	}
	return false, str
}