package models

type PaymentInfo struct {
	Action        string  `json:"action"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	RecipientName string  `json:"recipient_name"`
	AccountNumber string  `json:"account_number"`
	SortCode      string  `json:"sort_code"`
}
