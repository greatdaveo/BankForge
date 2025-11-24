package services

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func SendWhatsAppMessage(to string, message string) error {
	accountSID := os.Getenv("TWILLIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILLIO_AUTH_TOKEN")
	from := os.Getenv("TWILLIO_PHONE_NUMBER")

	// fmt.Println(".env: ", accountSID, authToken, from)

	// To create the req payload
	form := url.Values{}
	form.Add("To", to)
	form.Add("From", from)
	form.Add("Body", message)

	req, err := http.NewRequest(
		"POST",
		"https://api.twilio.com/2010-04-01/Accounts/"+accountSID+"/Messages.json",
		strings.NewReader(form.Encode()),
	)

	if err != nil {
		return err
	}

	req.SetBasicAuth(accountSID, authToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// To send the request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode >= 200 && response.StatusCode < 300 {
		return nil
	}

	return fmt.Errorf("failed to send message: status %d", response.StatusCode)
}
