package repo

import (
	"database/sql"
	"restApp/assets"
)

// CustomerRepository represents a data repository and implements CRUD methods for customers.
type CustomerRepository struct {
	db *sql.DB
}

// GetCustomerByID returns a single customer under the specified ID.
func (repo *CustomerRepository) GetCustomerByID(id int64) (*Customer, error) {
	script, err := assets.Asset("sql/customers/get_customer_by_id.sql")

	if err != nil {
		return nil, err
	}

	row := repo.db.QueryRow(string(script), id)
	customer := new(Customer)
	err = row.Scan(&customer.ID, &customer.Name, &customer.Address,
		&customer.TaxID, &customer.Email, &customer.PhoneNumber)

	if err != nil {
		return nil, err
	}

	return customer, nil
}

// GetAllCustomers returns a set of all customers from the database.
func (repo *CustomerRepository) GetAllCustomers() ([]*Customer, error) {
	script, err := assets.Asset("sql/customers/get_all_customers.sql")

	if err != nil {
		return nil, err
	}

	rows, err := repo.db.Query(string(script))

	if err != nil {
		return nil, err
	}

	customers := make([]*Customer, 0)

	for rows.Next() {
		customer := new(Customer)

		err = rows.Scan(&customer.ID, &customer.Name, &customer.Address,
			&customer.TaxID, &customer.Email, &customer.PhoneNumber)

		if err != nil {
			return nil, err
		}

		customers = append(customers, customer)
	}

	return customers, nil
}

// AddCustomer adds a new customer to the database.
func (repo *CustomerRepository) AddCustomer(customer *Customer) error {
	script, err := assets.Asset("sql/customers/add_customer.sql")

	if err != nil {
		return err
	}

	_, err = repo.db.Exec(string(script), customer.Name, customer.Address,
		customer.TaxID, customer.Email, customer.PhoneNumber)

	return err
}

// UpdateCustomer updates the customer in the database.
func (repo *CustomerRepository) UpdateCustomer(customer *Customer) error {
	script, err := assets.Asset("sql/customers/update_customer.sql")

	if err != nil {
		return err
	}

	_, err = repo.db.Exec(string(script), customer.ID, customer.Name, customer.Address,
		customer.TaxID, customer.Email, customer.PhoneNumber)

	return err
}

// DeleteCustomer deletes the customer from the database.
func (repo *CustomerRepository) DeleteCustomer(id int64) error {
	script, err := assets.Asset("sql/customers/delete_customer.sql")

	if err != nil {
		return err
	}

	_, err = repo.db.Exec(string(script), id)

	return err
}

// NewCustomerRepo creates a new repository for customers.
func NewCustomerRepo(db *sql.DB) *CustomerRepository {
	return &CustomerRepository{db}
}
