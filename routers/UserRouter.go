package routers

import (
	"net/http"
	"restTestOne/controllers"
)
var users controllers.UserController
// main routes function
func UserRouters() {
	http.HandleFunc("/user", userRouter)
	http.HandleFunc("/user/signin", signIn)
}
// main routes
func userRouter(w http.ResponseWriter, r *http.Request)  {
	query := r.URL.Query()["id"]

	switch {
		case r.Method == "GET" && len(query) > 0:
			users.GetOne(w, r)
		case r.Method == "GET":
			users.GetAll(w, r)
		// routes for signup
		case r.Method == "POST":
			users.Create(w, r)
	}
}
// post route for signin
func signIn(w http.ResponseWriter, r *http.Request) {
	users.GetOne(w, r)
}