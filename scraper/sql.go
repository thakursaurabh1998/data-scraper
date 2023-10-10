package scraper

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
)

var db *sql.DB

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "scraped_order_data"
)

const insertOrderQuery = `
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

const insertHistoricalQuery = `
INSERT INTO historical_data (id, remarks, awb_number, created_at, status, order_number)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (id) DO UPDATE SET
remarks = excluded.remarks,
awb_number = excluded.awb_number,
created_at = excluded.created_at,
status = excluded.status,
order_number = excluded.order_number`

func InsertOrderData(orderData *OrderData, awbNumber string, clientName string) {
	hashedId := generateHash([]string{orderData.OrderNumber, awbNumber})

	converted_amount, err := strconv.ParseFloat(orderData.Amount, 64)

	if err != nil {
		converted_amount = 0.0
	}

	executeQuery(insertOrderQuery, hashedId, orderData.CustomerName, orderData.CustomerDeliveryAddress, orderData.OrderNumber, orderData.CreatedAt, converted_amount, orderData.PaymentType, orderData.Client, awbNumber, clientName)
}

func InsertHistoricalData(historicalData *HistoricalData, order_number string, index int) {
	hashedId := generateHash([]string{historicalData.AWBNumber, order_number, strconv.Itoa(index)})

	executeQuery(insertHistoricalQuery, hashedId, historicalData.Remarks, historicalData.AWBNumber, historicalData.CreatedAt, historicalData.Status, order_number)
}

func executeQuery(query string, args ...any) {
	dbConn := GetDBConnection()
	_, err := dbConn.Exec(query, args...)
	if err != nil {
		log.Println("Error in query execution: ", err)
	}
}

func GetDBConnection() *sql.DB {
	var err error
	if db == nil {
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)
		db, err = sql.Open("postgres", psqlInfo)
		if err != nil {
			panic(err)
		}
		err = db.Ping()
		if err != nil {
			panic(err)
		}
	}

	return db
}
