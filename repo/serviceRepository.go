package repo

import "database/sql"

// ServiceRepository represents a data repository and implements CRUD methods for services.
type ServiceRepository struct {
	db *sql.DB
}

// GetServiceByID returns a single service under the specified ID.
func (repo *ServiceRepository) GetServiceByID(id int64) (*Service, error) {
	script, err := GetSQLScript("sql/services/get_service_by_id.sql")

	if err != nil {
		return nil, err
	}

	row := repo.db.QueryRow(script, id)
	service := new(Service)
	err = row.Scan(&service.ID, &service.Title, &service.Description, &service.Price)

	if err != nil {
		return nil, err
	}

	return service, nil
}

// GetAllServices returns a set of all services from the database.
func (repo *ServiceRepository) GetAllServices() ([]*Service, error) {
	script, err := GetSQLScript("sql/services/get_all_services.sql")

	if err != nil {
		return nil, err
	}

	rows, err := repo.db.Query(script)

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

// AddService adds a new service to the database.
func (repo *ServiceRepository) AddService(service *Service) error {
	script, err := GetSQLScript("sql/services/add_service.sql")

	if err != nil {
		return err
	}

	_, err = repo.db.Exec(script, service.Title, service.Description, service.Price)

	return err
}

// UpdateService updates the service in the database.
func (repo *ServiceRepository) UpdateService(service *Service) error {
	script, err := GetSQLScript("sql/services/update_service.sql")

	if err != nil {
		return err
	}

	_, err = repo.db.Exec(script, service.ID, service.Title,
		service.Description, service.Price)

	return err
}

// DeleteService deletes the service from the database.
func (repo *ServiceRepository) DeleteService(id int64) error {
	script, err := GetSQLScript("sql/services/delete_service.sql")

	if err != nil {
		return err
	}

	_, err = repo.db.Exec(script, id)

	return err
}

// NewServiceRepo creates a new repository for services.
func NewServiceRepo(db *sql.DB) *ServiceRepository {
	return &ServiceRepository{db}
}
