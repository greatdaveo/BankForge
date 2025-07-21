package controllers

import (
	"fmt"
	"net/http"

	"github.com/greatdaveo/SendlyPay/services"
)

func HandleIncomingMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// To parse form data sent by Twillio
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// To extract whatsapp messages and sender
	message := r.FormValue("Body")
	from := r.FormValue("From")

	fmt.Printf("Message from %s: %s\n", from, message)

	// To Call AI to extract payment details
	info, err := services.ExtractPaymentInfo(message)
	if err != nil {
		fmt.Printf("Failed to extract payment info: %v\n", err)
		http.Error(w, "Failed to process message", http.StatusInternalServerError)
		return
	}

	fmt.Println("Extracted Payment Info: -----------------------")
	fmt.Printf("Action: %s\n", info.Action)
	fmt.Printf("Amount: %.2f %s\n", info.Amount, info.Currency)
	fmt.Printf("Recipient: %s\n", info.RecipientName)
	fmt.Printf("Account Number: %s\n", info.AccountNumber)
	fmt.Printf("Sort Code: %s\n", info.SortCode)

	w.WriteHeader(http.StatusOK)
}
