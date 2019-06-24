package repo

// ICustomerRepository provides CRUD interface for customers.
type ICustomerRepository interface {
	GetCustomerByID(id int64) (*Customer, error)
	GetAllCustomers() ([]*Customer, error)
	AddCustomer(customer *Customer) error
	UpdateCustomer(customer *Customer) error
	DeleteCustomer(id int64) error
}

// IServiceRepository provides CRUD interface for services.
type IServiceRepository interface {
	GetServiceByID(id int64) (*Service, error)
	GetAllServices() ([]*Service, error)
	AddService(service *Service) error
	UpdateService(service *Service) error
	DeleteService(id int64) error
}

// IOrderRepository provides CRUD interface for orders.
type IOrderRepository interface {
	GetOrderByID(id int64) (*Order, error)
	GetAllOrders() ([]*Order, error)
	AddOrder(order *Order) error
	UpdateOrder(order *Order) error
	DeleteOrder(id int64) error
	GetOrderServiceByID(orderID int64, serviceID int64) (*Service, error)
	GetAllOrderServices(orderID int64) ([]*Service, error)
	AddServiceToOrder(orderID int64, serviceID int64) error
	DeleteServiceFromOrder(orderID int64, serviceID int64) error
}
