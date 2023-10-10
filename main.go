package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

const (
	URL      = "xxxx"
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "scraped_order_data"
	MAX_CONC = 100
)

func callAPI(generatedAWB string) (error, *Response) {
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

func fetchOrderForAWB(db *sql.DB, AWBNumber string) {
	// Call the API
	err, response := callAPI(AWBNumber)
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

	InsertOrderData(db, &response.Data.OrderData, AWBNumber, response.Data.Client)
	for index, historical_data := range response.Data.HistoricalData {
		InsertHistoricalData(db, &historical_data, response.Data.OrderData.OrderNumber, index)
	}

	log.Printf("Order inserted for AWB: %s", AWBNumber)
}

func DBConnection() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

func main() {
	db := DBConnection()
	defer db.Close()

	// SAMPLE AWB PRD000085800
	// start := 85635
	start := 150000
	end := 200000

	sem := make(chan int, MAX_CONC)

	for generatedAWB := start; generatedAWB <= end; generatedAWB++ {
		AWBNumber := fmt.Sprintf("PRD%09d", generatedAWB)

		log.Println("AWB:", AWBNumber)

		sem <- 1
		go func() {
			fetchOrderForAWB(db, AWBNumber)
			<-sem
		}()
	}
}

func InsertOrderData(db *sql.DB, orderData *OrderData, awbNumber string, clientName string) {
	hashedId := generateHash([]string{orderData.OrderNumber, awbNumber})
	sqlStatement := `
INSERT INTO orders (id, customer_name, customer_delivery_address, order_number, created_at, amount, payment_type, client, awb_number, client_name)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
ON CONFLICT (id) DO UPDATE SET
customer_name = excluded.customer_name,
customer_delivery_address = excluded.customer_delivery_address,
created_at = excluded.created_at,
amount = excluded.amount,
payment_type = excluded.payment_type,
client = excluded.client,
awb_number = excluded.awb_number,
client_name = excluded.client_name`

	converted_amount, err := strconv.ParseFloat(orderData.Amount, 64)

	if err != nil {
		converted_amount = 0.0
	}

	_, err = db.Exec(sqlStatement, hashedId, orderData.CustomerName, orderData.CustomerDeliveryAddress, orderData.OrderNumber, orderData.CreatedAt, converted_amount, orderData.PaymentType, orderData.Client, awbNumber, clientName)
	if err != nil {
		log.Println(err)
	}
}

func InsertHistoricalData(db *sql.DB, historicalData *HistoricalData, order_number string, index int) {
	hashedId := generateHash([]string{historicalData.AWBNumber, order_number, strconv.Itoa(index)})
	sqlStatement := `
INSERT INTO historical_data (id, remarks, awb_number, created_at, status, order_number)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (id) DO UPDATE SET
remarks = excluded.remarks,
awb_number = excluded.awb_number,
created_at = excluded.created_at,
status = excluded.status,
order_number = excluded.order_number`
	_, err := db.Exec(sqlStatement, hashedId, historicalData.Remarks, historicalData.AWBNumber, historicalData.CreatedAt, historicalData.Status, order_number)
	if err != nil {
		log.Println(err)
	}
}

func generateHash(arr []string) string {
	// Concatenate all strings in the array
	str := ""
	for _, s := range arr {
		str += s
	}

	// Generate hash using SHA256 algorithm
	hash := sha256.Sum256([]byte(str))

	// Convert hash to string and return
	return hex.EncodeToString(hash[:])
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
