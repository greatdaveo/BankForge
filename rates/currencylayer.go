package rates

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type CurrencyLayerAPI struct {
	APIKey string
}

func (c *CurrencyLayerAPI) FetchRate(from, to string) (float64, error) {
	from = strings.ToUpper(from)
	to = strings.ToUpper(to)

	url := fmt.Sprintf(
		"http://api.currencylayer.com/live?access_key=%s&source=%s&currencies=%s",
		c.APIKey, from, to,
	)
	println(url)

	response, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	var data struct {
		Success         bool               `json:"success"`
		ConversionRates map[string]float64 `json:"conversion_rates"`
	}

	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return 0, err
	}

	if !data.Success {
		return 0, errors.New("CurrencyLayer: API call failed")
	}

	key := from + to

	rate, ok := data.ConversionRates[key]
	if !ok {
		return 0, errors.New("CurrencyLayer: rate not found for pair " + key)
	}

	fmt.Println("CurrencyLayer rate: ", rate)
	return rate, nil

}
