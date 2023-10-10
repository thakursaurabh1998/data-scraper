package scraper

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

const URL = "xxxxx"

func Request(generatedAWB string) (error, *Response) {
	form := url.Values{}
	form.Add("awbnumber", generatedAWB)

	// Create a new HTTP request
	req, err := http.NewRequest(http.MethodPost, URL, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		return err, nil
	}

	// Send the request using the default HTTP client
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err, nil
	}
	defer resp.Body.Close()

	// Parse the response body as JSON
	var data Response
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return err, nil
	}

	return nil, &data
}

type OrderData struct {
	CustomerName            string `json:"CustomerName"`
	CustomerDeliveryAddress string `json:"CustomerDeliveryAddress"`
	OrderNumber             string `json:"OrderNumber"`
	CreatedAt               string `json:"CreatedAt"`
	Amount                  string `json:"Amount"`
	PaymentType             string `json:"PaymentType"`
	Client                  int    `json:"Client"`
}

type HistoricalData struct {
	Remarks   string `json:"Remarks"`
	AWBNumber string `json:"AWBNumber"`
	CreatedAt string `json:"CreatedAt"`
	Status    string `json:"Status"`
}

type ResponseData struct {
	OrderData      OrderData        `json:"orderData"`
	HistoricalData []HistoricalData `json:"histData"`
	Client         string           `json:"Client"`
}

type Response struct {
	Status     bool         `json:"Status"`
	StatusCode int16        `json:"StatusCode"`
	Message    string       `json:"message"`
	Data       ResponseData `json:"data"`
}
