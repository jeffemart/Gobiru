package main

func main() {
	router := SetupRouter()
	router.Run(":8080")
}
