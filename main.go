package main

import (
	"fmt"
	"log"
	"net/http"
	"restTestOne/routers"
)

func main() {
	fmt.Println("start main")

	routers.UserRouters()
	routers.ProductRouters()

	log.Fatal(http.ListenAndServe(":4000", nil))
}