package repo

import (
	"time"
)

// Customer represents a single customer of the company.
type Customer struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	TaxID       string `json:"tax_id"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

// Service represents a single service provided by the company.
type Service struct {
	ID          int64   `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

// Order represents a single order made by some of the company's customers.
type Order struct {
	ID         int64     `json:"id"`
	CustomerID int64     `json:"customer_id"`
	Date       time.Time `json:"date"`
}
