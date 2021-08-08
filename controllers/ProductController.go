package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"restTestOne/models"
	"restTestOne/response"
	"restTestOne/services"
	"restTestOne/utils"
	"strconv"
	"strings"
)

type ProductController struct {}

var product = services.ProductService{}
var user = services.UserService{}
// get all product controller
func (p *ProductController) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	res, err := product.GetAll()
	if err != nil {
		res := response.ErrorResponse{Status: "failed to get data", Message: err.Error()}
		json.NewEncoder(w).Encode(res)
		return
	}

	json.NewEncoder(w).Encode(res)
}
// create or update controller
func (p *ProductController) CreateOrUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// authentication check
	auth := r.Header.Get("Authorization")
	id, _,errCheck := authCheck(auth)
	if errCheck != nil {
		fmt.Println(errCheck)
		json.NewEncoder(w).Encode(errCheck)
		return
	}
	var pr models.Product
	// handle x-www-form-urlencoded
	r.ParseMultipartForm(0)
	for k,v := range r.Form {
		if len(r.Form["id"][0]) > 0 {
			break
		}
		// check input form
		check, str := contains(k, []string{"productName", "productPrice", "productStock"})
		if !check {
			json.NewEncoder(w).Encode(str+" is not a valid input")
			return
		}
		// set form data to struct
		if k == "productName" {
			pr.ProductName = v[0]
		}
		if k == "productPrice" {
			intVar, _ := strconv.Atoi(v[0])
			pr.ProductPrice = intVar
		}
		if k == "productStock" {
			intVar, _ := strconv.Atoi(v[0])
			pr.ProductStock = intVar
		}
	}
	// parse JSON stringify or JSON object
	json.NewDecoder(r.Body).Decode(&pr)
	pr.SuplierId = id
	var err error

	if r.Method == "POST" {
		pr, err = product.Create(pr)
	} else if r.Method == "PUT"{
		pr, err = product.Update(pr, r.URL.Query()["id"][0])
	}
	
	if err != nil {
		res := response.ErrorResponse{Status: "failed create Product", Message: err.Error()}
		json.NewEncoder(w).Encode(res)
		return
	}

	response := response.CreatePrResponse{Status: "Success create Product", Data: pr}
	json.NewEncoder(w).Encode(response)
}
// buy or add controller
func (p *ProductController) Buy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// authentication check
	auth := r.Header.Get("Authorization")
	id, types, err := authCheck(auth)
	if err != nil {
		fmt.Println(err)
		json.NewEncoder(w).Encode(err)
		return
	}

	var pr models.EditProduct
	
	pr.CreatedBy = id
	pr.Type = types
	// handle x-www-form-urlencoded
	r.ParseMultipartForm(0)
	for k,v := range r.Form {
		if len(r.Form["id"][0]) > 0 {
			break
		}
		// check input form
		check, str := contains(k, []string{"productId", "amount"})
		if !check {
			json.NewEncoder(w).Encode(str+" is not a valid input")
			return
		}
		// set form data to struct
		if k == "productId" {
			pr.ProductId = v[0]
		}
		if k == "amount" {
			intVar, _ := strconv.Atoi(v[0])
			pr.Amount = intVar
		}
	}
	// parse JSON stringify or JSON object
	json.NewDecoder(r.Body).Decode(&pr)
	
	result, err := product.BuyOrAdd(pr)
	if err != nil {
		res := response.ErrorResponse{Status: "transaction buy failed", Message: err.Error()}
		json.NewEncoder(w).Encode(res)
		return
	}

	json.NewEncoder(w).Encode(result)
}
// auth check helper
func authCheck(auth string ) (string, string, error) {
	if len(auth) == 0 {
		return "", "", fmt.Errorf("token is required")
	}

	splitToken := strings.Split(auth, "Bearer ")
	auth = splitToken[1]
	status, token := utils.TokenDecrypt(auth)
	if !status {
		return "", "", fmt.Errorf("token not valid")
	}

	isAllow, result, types := user.CheckUser(token)
	if !isAllow {
		return "", "", fmt.Errorf(result)
	}

	return result, types, nil
}