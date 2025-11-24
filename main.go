package main

import (
	"log"
	"net/http"
	"os"

	"github.com/greatdaveo/XChange/controllers"
	"github.com/greatdaveo/XChange/rates"
	"github.com/greatdaveo/XChange/routes"
	"github.com/greatdaveo/XChange/services"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	apiKey := os.Getenv("EXCHANGE_RATE_API_KEY")
	currencyLayerKey := os.Getenv("CURRENCY_LAYER_API_KEY")

	externalAPI := &rates.ExchangeRateAPI{
		APIKey: apiKey,
	}
	frankfurterAPI := &rates.CurrencyLayerAPI{
		APIKey: currencyLayerKey,
	}

	//  Pass both as RateProvider interface
	rateService := services.NewRateService([]rates.RateProvider{
		externalAPI,
		frankfurterAPI,
	})

	r := routes.RegisterRoutes(controllers.ConvertCurrency(rateService))

	log.Println("📌 Server running on :8080")
	http.ListenAndServe(":8080", r)
}
