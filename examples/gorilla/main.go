package main

import (
	"log"
	"net/http"

	"github.com/jeffemart/Gobiru/examples/gorilla/complete/routes"
)

func main() {
	router := routes.SetupRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}
