package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/greatdaveo/SendlyPay/routes"
)

func main() {
	routes.WhatsAppRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("📌 SendlyPay is running on port!!! ", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("❌ Could not start sever: %v", err)
	}

}
