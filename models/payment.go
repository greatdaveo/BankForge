package models

type PaymentInfo struct {
	Action        string  `json:"action"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	RecipientName string  `json:"recipient_name"`
	AccountNumber string  `json:"account_number"`
	SortCode      string  `json:"sort_code"`
}

type TransactionReceipt struct {
	Reference     string
	FromUser      string
	ToName        string
	AccountNumber string
	SortCode      string
	Amount        float64
	Currency      string
	Status        string
}

// For BuildPaymentPayload func

type PaymentPayload struct {
	User          string        `json:"user"`
	AmountInMinor int           `json:"account_in_minor"`
	Currency      string        `json:"currency"`
	PaymentMethod PaymentMethod `json:"payment_method"`
}
type PaymentMethod struct {
	Type              string            `json:"type"`
	ProviderSelection ProviderSelection `json:"provider_selection"`
	Beneficiary       Beneficiary       `json:"beneficiary"`
}

type ProviderSelection struct {
	Type string `json:"type"`
}

type Beneficiary struct {
	Type              string            `json:"type"`
	AccountHolderName string            `json:"account_holder_name"`
	Reference         string            `json:"reference"`
	AccountIdentifier AccountIdentifier `json:"account_identifier"`
}

type AccountIdentifier struct {
	Type          string `json:"type"`
	AccountNumber string `json:"account_number"`
	SortCode      string `json:"sort_code"`
}
