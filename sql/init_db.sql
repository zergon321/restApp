CREATE TABLE customers (
    id SERIAL PRIMARY KEY,
    company_name VARCHAR (128) NOT NULL,
    company_address VARCHAR (256) NOT NULL,
    tax_id CHAR (12) UNIQUE NOT NULL,
    email VARCHAR (256) UNIQUE NOT NULL,
    phone_number CHAR(14) UNIQUE NOT NULL
);

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    customer_id SERIAL REFERENCES customers ON DELETE CASCADE,
    contract_date DATE NOT NULL DEFAULT CURRENT_DATE
);

CREATE TABLE services (
    id SERIAL PRIMARY KEY,
    title VARCHAR (256) UNIQUE NOT NULL,
    service_description VARCHAR (512) NOT NULL,
    price DECIMAL (9, 2) NOT NULL 
);

CREATE TABLE orders_to_services (
    order_id SERIAL REFERENCES orders,
    service_id SERIAL REFERENCES services,
    PRIMARY KEY (order_id, service_id)
);