package main

import (
	"log"
	"net/http"

	"github.com/jeffemart/gobiru/examples/gorilla/routes"
)

func main() {
	r := routes.SetupRoutes()
	log.Fatal(http.ListenAndServe(":3000", r))
}
