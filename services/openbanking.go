package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	tlsigning "github.com/Truelayer/truelayer-signing/go"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/greatdaveo/SendlyPay/models"
)

type OAuthToken struct {
	AccessToken string `json:"access_token"`
	User        string `json:"user"`
}

func ExchangeCodeForToken(code string) (OAuthToken, error) {
	// To prepare req body
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("client_id", os.Getenv("TRUELAYER_CLIENT_ID"))
	form.Set("client_secret", os.Getenv("TRUELAYER_CLIENT_SECRET"))
	form.Set("redirect_uri", os.Getenv("TRUELAYER_REDIRECT_URI"))
	form.Set("code", code)

	req, err := http.NewRequest("POST", os.Getenv("TRUELAYER_AUTH_URL")+"/connect/token", bytes.NewBufferString(form.Encode()))
	if err != nil {
		return OAuthToken{}, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// To send request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return OAuthToken{}, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	// fmt.Println("Raw token response:", string(body))

	if err != nil {
		return OAuthToken{}, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return OAuthToken{}, err
	}

	fmt.Println(" Keys in result:", result)

	token, _ := result["access_token"].(string)
	// user, _ := result["user"].(string)

	var user string

	parts := strings.Split(token, ".")
	if len(parts) == 3 {
		claims := jwt.MapClaims{}
		_, _, err := new(jwt.Parser).ParseUnverified(token, claims)
		if err == nil {
			if sub, ok := claims["sub"].(string); ok {
				user = sub
			}
		}
	}

	if token == "" {
		return OAuthToken{}, fmt.Errorf("access_token is missing")
	}
	if user == "" {
		return OAuthToken{}, fmt.Errorf("user field is missing")
	}

	return OAuthToken{
			AccessToken: token,
			User:        user,
		},
		nil
}

func InitiatePayment(token string, paymentPayload models.PaymentPayload) error {

	// To load the private key
	privateKeyBytes, err := os.ReadFile("ec512-private-key.pem")
	if err != nil {
		return fmt.Errorf("failed to load private key: %v", err)
	}

	// To generate a fresh UUID for idempotency key
	idempotencyKey := uuid.New().String()

	paymentBody, err := json.Marshal(paymentPayload)
	if err != nil {
		return fmt.Errorf("failed to marshal payment body: %v", err)
	}

	fmt.Println("Final JSON:", string(paymentBody))
	fmt.Println("Body length:", len(paymentBody))

	// To sign the request
	signature, err := tlsigning.SignWithPem(os.Getenv("TRUELAYER_KID"), privateKeyBytes).
		Method("POST").
		Path("/v3/payments").
		Header("Idempotency-Key", []byte(idempotencyKey)).
		Body(paymentBody).
		Sign()

	if err != nil {
		return fmt.Errorf("failed to sign request: %v", err)
	}

	// To send the request to True layer
	req, err := http.NewRequest("POST", os.Getenv("TRUELAYER_API_URL")+"/v3/payments", bytes.NewBuffer(paymentBody))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotency-Key", idempotencyKey)
	req.Header.Set("Tl-Signature", signature)

	// fmt.Println("Signed Idempotency-Key:", string([]byte(idempotencyKey)))
	// fmt.Println("HTTP Idempotency-Key:", req.Header.Get("Idempotency-Key"))

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	fmt.Println("Status Code: ", response.StatusCode)
	bodyResponse, _ := io.ReadAll(response.Body)
	fmt.Println("Response body: ", string(bodyResponse))

	if response.StatusCode >= 200 && response.StatusCode < 300 {
		return nil
	}

	return fmt.Errorf("payment failed: %s", string(bodyResponse))
}

func BuildPaymentPayload(info models.PaymentInfo, userID string) (models.PaymentPayload, error) {
	sortCode := strings.ReplaceAll(info.SortCode, "-", "")
	sortCode = strings.ReplaceAll(sortCode, " ", "")

	payload := models.PaymentPayload{
		User:          userID,
		AmountInMinor: int(info.Amount * 100),
		Currency:      info.Currency,
		PaymentMethod: models.PaymentMethod{
			Type: "bank_transfer",
			ProviderSelection: models.ProviderSelection{
				Type: "user_selected",
			},
			Beneficiary: models.Beneficiary{
				Type:              "external_account",
				AccountHolderName: info.RecipientName,
				Reference:         "SendlyPay transfer",
				AccountIdentifier: models.AccountIdentifier{
					Type:          "sort_code_account_number",
					AccountNumber: info.AccountNumber,
					SortCode:      sortCode,
				},
			},
		},
	}

	return payload, nil
}
