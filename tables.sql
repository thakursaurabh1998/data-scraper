CREATE TABLE orders (
	id VARCHAR(255) PRIMARY KEY,
	customer_name VARCHAR(255),
	customer_delivery_address VARCHAR(2000),
	awb_number VARCHAR(255),
	order_number VARCHAR(255),
	created_at TIMESTAMP,
	amount VARCHAR(255),
	payment_type VARCHAR(255),
	client INT
	client_name VARCHAR(255)
);

CREATE TABLE historical_data (
	id VARCHAR(255) PRIMARY KEY,
	remarks VARCHAR(255),
	order_number VARCHAR(255),
	awb_number VARCHAR(255),
	created_at TIMESTAMP,
	status VARCHAR(255)
);
