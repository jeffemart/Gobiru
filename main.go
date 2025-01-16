package main

import (
	"log"
)

func main() {
	router := SetupRouter()

	// Servir a documentação
	router.Static("/docs", "./docs")

	log.Fatal(router.Run(":8080"))
}
