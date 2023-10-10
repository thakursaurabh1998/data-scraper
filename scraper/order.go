package scraper

import "log"

func FetchOrderForAWB(AWBNumber string) {
	log.Printf("Fetching order for AWB: %s", AWBNumber)
	// Call the API
	err, response := Request(AWBNumber)
	if err != nil || response == nil {
		log.Print(err)
		return
	}

	if response.StatusCode != 200 {
		log.Printf("Order not found for AWB: %s", AWBNumber)
		return
	}

	// Print the response
	log.Printf("%+v\n", response)

	InsertOrderData(&response.Data.OrderData, AWBNumber, response.Data.Client)
	for index, historical_data := range response.Data.HistoricalData {
		InsertHistoricalData(&historical_data, response.Data.OrderData.OrderNumber, index)
	}

	log.Printf("Order inserted for AWB: %s", AWBNumber)
}
