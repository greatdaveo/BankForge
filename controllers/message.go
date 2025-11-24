package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/greatdaveo/SendlyPay/models"
	"github.com/greatdaveo/SendlyPay/services"
	"github.com/greatdaveo/SendlyPay/store"
)

// In memory session
var sessionStore = make(map[string]models.PaymentInfo)

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
	message := strings.TrimSpace(strings.ToLower(r.FormValue("Body")))
	from := r.FormValue("From")
	fmt.Printf("📩 Message from %s: %s\n", from, message)

	// To detect the "my transactions" for getting past transactions
	if message == "my transactions" || message == "my transaction" {
		totalTransactions := services.GetRecentTransactions(from, 5)

		if len(totalTransactions) == 0 {
<<<<<<< HEAD
			_ = services.SendWhatsAppMessage(from, "No transactions found.")
=======
			_ = services.SendWhatsAppMessage(from, " No transactions found.")
>>>>>>> 5928aff (Removed unnessary details from logs)
		} else {
			reply := "📑 Your recent transactions:\n"
			for _, transaction := range totalTransactions {
				line := fmt.Sprintf("- %s: £%.2f to %s (%s)\n", transaction.Reference, transaction.Amount, transaction.ToName, transaction.Status)
				reply += line
			}
			_ = services.SendWhatsAppMessage(from, reply)
		}

		w.WriteHeader(http.StatusOK)
		return
	}

	// To check if user has a pending payment in session
	pending, hasSession := sessionStore[from]

	// If yes confirm and process payment
	if hasSession && message == "yes" {
		reply := fmt.Sprintf("🔄 Payment of £%.2f to %s is being processed.", pending.Amount, pending.RecipientName)
		_ = services.SendWhatsAppMessage(from, reply)

		if services.DeductFromWallet(from, pending.Amount) {
			receipt := services.GenerateReceipt(from, pending)

			// To retrieve the access token
			oauth, ok := store.AccessTokenStore[from]
			if !ok {
				services.SendWhatsAppMessage(from, "You need to link your bank before making payments.")
				return
			}

			paymentPayload, err := services.BuildPaymentPayload(pending, oauth.User)
			if err != nil {
				fmt.Println("Failed to build payment body: ", err)
				services.SendWhatsAppMessage(from, "Failed to process payment payload.")
				return
			}

			err = services.InitiatePayment(oauth.AccessToken, paymentPayload)
			if err != nil {
				fmt.Println("Payment request to TrueLayer failed:", err)
				services.SendWhatsAppMessage(from, "Payment failed to process.")
				return
			}

			services.SendWhatsAppMessage(from, " Payment successfully initiated with the bank.")

			// To save receipt to the user history
			services.SaveTransaction(from, receipt)

			// To send successful message
			reply := fmt.Sprintf(
<<<<<<< HEAD
				" Payment Successful!\nRef: %s\nSent £%.2f to %s\nAccount Number: %s\nSort Code: %s",
=======
				"Payment Successful!\nRef: %s\nSent £%.2f to %s\nAccount Number: %s\nSort Code: %s",
>>>>>>> 5928aff (Removed unnessary details from logs)
				receipt.Reference,
				receipt.Amount,
				receipt.ToName,
				receipt.AccountNumber,
				receipt.SortCode,
			)
			_ = services.SendWhatsAppMessage(from, reply)
<<<<<<< HEAD
			fmt.Printf(" Payment processed for %s: £%.2f\n", from, pending.Amount)
		} else {
			reply := fmt.Sprintf("Insufficient balance. Your wallet has £%.2f", services.GetBalance(from))
			_ = services.SendWhatsAppMessage(from, reply)
			fmt.Printf("Payment failed for %s: insufficient funds\n", from)
=======
			fmt.Printf("Payment processed for %s: £%.2f\n", from, pending.Amount)
		} else {
			reply := fmt.Sprintf(" Insufficient balance. Your wallet has £%.2f", services.GetBalance(from))
			_ = services.SendWhatsAppMessage(from, reply)
			fmt.Printf(" Payment failed for %s: insufficient funds\n", from)
>>>>>>> 5928aff (Removed unnessary details from logs)
		}

		// To clear session
		delete(sessionStore, from)
		w.WriteHeader(http.StatusOK)
		return
	}

	// If No, cancel payment
	if hasSession && message == "no" {
<<<<<<< HEAD
		fmt.Printf("Payment canceled by user: %+v\n", pending)

		// To send cancellation reply
		_ = services.SendWhatsAppMessage(from, "Payment cancelled.")
=======
		fmt.Printf(" Payment canceled by user: %+v\n", pending)

		// To send cancellation reply
		_ = services.SendWhatsAppMessage(from, " Payment cancelled.")
>>>>>>> 5928aff (Removed unnessary details from logs)
		delete(sessionStore, from)
		w.WriteHeader(http.StatusOK)
		return
	}

	// To Call AI to extract payment details as a new message, if otherwise
	info, err := services.ExtractPaymentInfo(r.FormValue("Body"))
	if err != nil {
<<<<<<< HEAD
		fmt.Printf("Failed to extract payment info: %v\n", err)
=======
		fmt.Printf(" Failed to extract payment info: %v\n", err)
>>>>>>> 5928aff (Removed unnessary details from logs)
		http.Error(w, "Failed to process message", http.StatusInternalServerError)
		return
	}

	// To store in session for confirmation
	sessionStore[from] = info

	// To handle confirmation message
	reply := fmt.Sprintf(
<<<<<<< HEAD
		" Got it. You want to send £%.2f to %s (Account: %s, Sort Code: %s).\nReply with 'Yes' to confirm or 'No' to cancel.",
=======
		"Got it. You want to send £%.2f to %s (Account: %s, Sort Code: %s).\nReply with 'Yes' to confirm or 'No' to cancel.",
>>>>>>> 5928aff (Removed unnessary details from logs)
		info.Amount,
		info.RecipientName,
		info.AccountNumber,
		info.SortCode,
	)

	// To send confirmation message from WhatsApp
	if err := services.SendWhatsAppMessage(from, reply); err != nil {
<<<<<<< HEAD
		fmt.Printf("Failed to send confirmation: %v\n", err)
=======
		fmt.Printf(" Failed to send confirmation: %v\n", err)
>>>>>>> 5928aff (Removed unnessary details from logs)
		http.Error(w, "Failed to send message", http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusOK)
}
