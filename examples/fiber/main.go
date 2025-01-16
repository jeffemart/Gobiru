package main

func main() {
	app := SetupRouter()
	app.Listen(":8080")
}
