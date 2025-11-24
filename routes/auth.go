package routes

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/greatdaveo/SendlyPay/services"
	"github.com/greatdaveo/SendlyPay/store"
)

func AuthFlow(w http.ResponseWriter, r *http.Request) {
	clientID := os.Getenv("TRUELAYER_CLIENT_ID")
	redirectURI := os.Getenv("TRUELAYER_REDIRECT_URI")
	authBaseURL := os.Getenv("TRUELAYER_AUTH_URL")

	// fmt.Println(clientID, redirectURI, authBaseURL)

	query := url.Values{}
	query.Set("response_type", "code")
	query.Set("client_id", clientID)
	query.Set("redirect_uri", redirectURI)
	query.Set("scope", "payments")
	query.Set("providers", "uk-cs-mock")
	query.Set("state", "sendlypay123")

	authURL := fmt.Sprintf("%s/?%s", authBaseURL, query.Encode())
	// fmt.Println("Redirecting to:", authURL)
	http.Redirect(w, r, authURL, http.StatusFound)
}

func RegisterRoute() {
	http.HandleFunc("GET /auth/initiate", AuthFlow)
	http.HandleFunc("/auth/callback", AuthCallback)
}

func AuthCallback(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("Raw callback URL:", r.URL)
	// fmt.Println("Raw callback query:", r.URL.RawQuery)

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "No code received", http.StatusBadRequest)
		return
	}

	// To exchange the code for an access token
	oauth, err := services.ExchangeCodeForToken(code)
	if err != nil {
		fmt.Printf("Failed to exchange code: %v\n", err)
		http.Error(w, "Failed to get access token", http.StatusInternalServerError)
		return
	}

	store.AccessTokenStore["whatsapp:+447778797699"] = oauth

	// fmt.Println("Access Token: ", token)
	w.Write([]byte("Bank account linked successfully. You can return to WhatsApp."))
}
