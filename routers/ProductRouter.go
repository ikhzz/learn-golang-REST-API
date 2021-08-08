package routers

import (
	"encoding/json"
	"net/http"
	"restTestOne/controllers"
)

var product controllers.ProductController
// main routes function
func ProductRouters() {
	http.HandleFunc("/product", mainRoute)
	http.HandleFunc("/product/buy", buy)
}
// main routes
func mainRoute(w http.ResponseWriter, r *http.Request){
	query := r.URL.Query()["id"]
	switch {
		case r.Method == "GET":
			product.GetAll(w, r)
		case r.Method == "POST" || (len(query) > 0 && r.Method == "PUT"):
			product.CreateOrUpdate(w, r)
	}
}
// buy or add routes
func buy(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()["id"]
	if len(query) > 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("id product is required")
		return
	}
	product.Buy(w, r)
}