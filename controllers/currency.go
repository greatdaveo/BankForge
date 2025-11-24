package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/greatdaveo/XChange/services"
)

func ConvertCurrency(service *services.RateService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		query := r.URL.Query()
		amountStr := query.Get("amount")
		from := query.Get("from")
		to := query.Get("to")

		if amountStr == "" || from == "" || to == "" {
			http.Error(w, "Missing required query parameters", http.StatusBadRequest)
			return
		}

		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil || amount <= 0 {
			http.Error(w, "Invalid amount", http.StatusBadRequest)
			return
		}

		// To get real exchange rate from the RateService
		rate, err := service.GetRate(from, to)
		// fmt.Printf(" Received rate: %f\n", rate)
		if err != nil || rate <= 0 {
			fmt.Println(err)
			http.Error(w, "Unable to convert currency", http.StatusInternalServerError)
			return
		}

		convertedAmount := amount * rate

		result := map[string]interface{}{
			"from":             from,
			"to":               to,
			"original_amount":  amount,
			"converted_amount": convertedAmount,
			"rate_used":        rate,
			"retrieved_at":     time.Now(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}

}
