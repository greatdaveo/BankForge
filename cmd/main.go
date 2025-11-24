package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/greatdaveo/SendlyPay/routes"
	"github.com/greatdaveo/SendlyPay/services"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	routes.WhatsAppRoutes()
	routes.RegisterRoute()

	// To set a temporary initial amount balance
	services.SetInitialBalance("whatsapp:+447778797699", 200.00)

	port := os.Getenv("PORT")
	fmt.Println("PORT: ", port)
	if port == "" {
		port = "8080"
	}

	fmt.Println("📌 SendlyPay is running on port!!! ", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
<<<<<<< HEAD
		log.Fatalf("Could not start sever: %v", err)
=======
		log.Fatalf(" Could not start sever: %v", err)
>>>>>>> 5928aff (Removed unnessary details from logs)
	}

}
