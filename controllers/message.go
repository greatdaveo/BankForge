package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/greatdaveo/SendlyPay/models"
	"github.com/greatdaveo/SendlyPay/services"
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
			_ = services.SendWhatsAppMessage(from, "❌ No transactions found.")
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
			// To save receipt to the user history
			services.SaveTransaction(from, receipt)

			// To send successful message
			reply := fmt.Sprintf(
				"✅ Payment Successful!\nRef: %s\nSent £%.2f to %s\nAccount Number: %s\nSort Code: %s",
				receipt.Reference,
				receipt.Amount,
				receipt.ToName,
				receipt.AccountNumber,
				receipt.SortCode,
			)
			_ = services.SendWhatsAppMessage(from, reply)
			fmt.Printf("✅ Payment processed for %s: £%.2f\n", from, pending.Amount)
		} else {
			reply := fmt.Sprintf("❌ Insufficient balance. Your wallet has £%.2f", services.GetBalance(from))
			_ = services.SendWhatsAppMessage(from, reply)
			fmt.Printf("❌ Payment failed for %s: insufficient funds\n", from)
		}

		// To clear session
		delete(sessionStore, from)
		w.WriteHeader(http.StatusOK)
		return
	}

	// If No, cancel payment
	if hasSession && message == "no" {
		fmt.Printf("❌ Payment canceled by user: %+v\n", pending)

		// To send cancellation reply
		_ = services.SendWhatsAppMessage(from, "❌ Payment cancelled.")
		delete(sessionStore, from)
		w.WriteHeader(http.StatusOK)
		return
	}

	// To Call AI to extract payment details as a new message, if otherwise
	info, err := services.ExtractPaymentInfo(r.FormValue("Body"))
	if err != nil {
		fmt.Printf("❌ Failed to extract payment info: %v\n", err)
		http.Error(w, "Failed to process message", http.StatusInternalServerError)
		return
	}

	// To store in session for confirmation
	sessionStore[from] = info

	// To handle confirmation message
	reply := fmt.Sprintf(
		"✅ Got it. You want to send £%.2f to %s (Account: %s, Sort Code: %s).\nReply with 'Yes' to confirm or 'No' to cancel.",
		info.Amount,
		info.RecipientName,
		info.AccountNumber,
		info.SortCode,
	)

	// To send confirmation message from WhatsApp
	if err := services.SendWhatsAppMessage(from, reply); err != nil {
		fmt.Printf("❌ Failed to send confirmation: %v\n", err)
		http.Error(w, "Failed to send message", http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusOK)
}
