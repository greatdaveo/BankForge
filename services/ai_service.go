package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/greatdaveo/SendlyPay/models"
)

func ExtractPaymentInfo(message string) (models.PaymentInfo, error) {
	var info models.PaymentInfo

	systemPrompt := "You are a helpful assistant that extracts payment details from a short user message. Always respond ONLY with valid JSON."
	userPrompt := fmt.Sprintf(`Extract payment info from this message: "%s"

	Respond only in this format:

	{
	"action": "send",
	"amount": 100,
	"currency": "GBP",
	"recipient_name": "John Doe",
	"account_number": "12345678",
	"sort_code": "04-03-11"
	}
	`, message)

	reqBody := map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": userPrompt},
		},
		"temperature": 0,
	}

	// To marshall req payload to JSON
	body, err := json.Marshal(reqBody)
	if err != nil {
		return info, nil
	}

	// To send HTTP Req to OpenAI
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(body))
	if err != nil {
		return info, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return info, err
	}
	defer response.Body.Close()

	// To Parse GPT Response
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return info, err
	}

	// fmt.Println("responseBody:", string(responseBody))

	var result map[string]interface{}
	if err := json.Unmarshal(responseBody, &result); err != nil {
		return info, err
	}

	// choices, ok := result["choices"][0]["message"]["content"]
	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return info, errors.New("no choices returned from OPENAI")
	}

	choice := choices[0].(map[string]interface{})
	messageObj := choice["message"].(map[string]interface{})
	content := messageObj["content"].(string)

	// To parse AI JSON into struct
	err = json.Unmarshal([]byte(content), &info)
	if err != nil {
		return info, fmt.Errorf("failed to parse AI JSON: %v", err)
	}

	return info, nil

}
