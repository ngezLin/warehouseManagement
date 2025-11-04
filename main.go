package main

import (
	"fmt"
	"log"
	"net/http"
	"warehousemanagement/config"
	"warehousemanagement/routers"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	config.ConnectDB()
	config.MigrateDB()
	r := routers.SetupRouter()

	fmt.Println("Server running at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
