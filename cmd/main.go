package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/greatdaveo/SendlyPay/routes"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	routes.WhatsAppRoutes()

	port := os.Getenv("PORT")
	fmt.Println("PORT: ", port)
	if port == "" {
		port = "8080"
	}

	fmt.Println("📌 SendlyPay is running on port!!! ", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("❌ Could not start sever: %v", err)
	}

}
