package repo

import "database/sql"

// OrderRepository represents a data repository and implements CRUD methods for orders.
type OrderRepository struct {
	db *sql.DB
}

// GetOrderByID returns a single order under the specified ID.
func (repo *OrderRepository) GetOrderByID(id int64) (*Order, error) {
	script, err := GetSQLScript("sql/orders/get_order_by_id.sql")

	if err != nil {
		return nil, err
	}

	row := repo.db.QueryRow(script, id)
	order := new(Order)
	err = row.Scan(&order.ID, &order.CustomerID, &order.Date)

	if err != nil {
		return nil, err
	}

	return order, nil
}

// GetAllOrders returns a set of all orders from the database.
func (repo *OrderRepository) GetAllOrders() ([]*Order, error) {
	script, err := GetSQLScript("sql/orders/get_all_orders.sql")

	if err != nil {
		return nil, err
	}

	rows, err := repo.db.Query(script)

	if err != nil {
		return nil, err
	}

	orders := make([]*Order, 0)

	for rows.Next() {
		order := new(Order)
		err = rows.Scan(&order.ID, &order.CustomerID, &order.Date)

		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

// AddOrder adds a new order to the database.
func (repo *OrderRepository) AddOrder(order *Order) error {
	script, err := GetSQLScript("sql/orders/add_order.sql")

	if err != nil {
		return err
	}

	_, err = repo.db.Exec(script, order.CustomerID, order.Date)

	return err
}

// UpdateOrder updates the order in the database.
func (repo *OrderRepository) UpdateOrder(order *Order) error {
	script, err := GetSQLScript("sql/orders/update_order.sql")

	if err != nil {
		return err
	}

	_, err = repo.db.Exec(script, order.ID, order.Date)

	return err
}

// DeleteOrder deletes the order from the database.
func (repo *OrderRepository) DeleteOrder(id int64) error {
	script, err := GetSQLScript("sql/orders/delete_order.sql")

	if err != nil {
		return err
	}

	_, err = repo.db.Exec(script, id)

	return err
}

// GetOrderServiceByID returns a single service included in the order by its ID.
func (repo *OrderRepository) GetOrderServiceByID(orderID int64, serviceID int64) (*Service, error) {
	script, err := GetSQLScript("sql/orders/get_order_service_by_id.sql")

	if err != nil {
		return nil, err
	}

	row := repo.db.QueryRow(script, orderID, serviceID)
	service := new(Service)
	err = row.Scan(&service.ID, &service.Title, &service.Description, &service.Price)

	if err != nil {
		return nil, err
	}

	return service, nil
}

// GetAllOrderServices returns all the services included in the order.
func (repo *OrderRepository) GetAllOrderServices(orderID int64) ([]*Service, error) {
	script, err := GetSQLScript("sql/orders/get_all_order_services.sql")

	if err != nil {
		return nil, err
	}

	rows, err := repo.db.Query(script, orderID)

	if err != nil {
		return nil, err
	}

	services := make([]*Service, 0)

	for rows.Next() {
		service := new(Service)

		err = rows.Scan(&service.ID, &service.Title, &service.Description, &service.Price)

		if err != nil {
			return nil, err
		}

		services = append(services, service)
	}

	return services, nil
}

// AddServiceToOrder adds a service to the order.
func (repo *OrderRepository) AddServiceToOrder(orderID int64, serviceID int64) error {
	script, err := GetSQLScript("sql/orders/add_service_to_order.sql")

	if err != nil {
		return err
	}

	_, err = repo.db.Exec(script, orderID, serviceID)

	return err
}

// DeleteServiceFromOrder deleted the service from the order.
func (repo *OrderRepository) DeleteServiceFromOrder(orderID int64, serviceID int64) error {
	script, err := GetSQLScript("sql/orders/delete_service_from_order.sql")

	if err != nil {
		return err
	}

	_, err = repo.db.Exec(script, orderID, serviceID)

	return err
}

// NewOrderRepository creates a new repository for orders and their services.
func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db}
}
