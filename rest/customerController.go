package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"restApp/repo"
	"strconv"

	"github.com/gorilla/mux"
)

// CustomerController provides REST API methods for customers.
type CustomerController struct {
	customerRepo repo.ICustomerRepository
	controller
}

func (ctl *CustomerController) getCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		ctl.handleWebError(w, http.StatusBadRequest,
			fmt.Sprintf("Incorrect parameter for id: %v", params["id"]))

		return
	}

	customer, err := ctl.customerRepo.GetCustomerByID(int64(id))

	if err != nil {
		ctl.handleInternalError("Database access error", err)
		ctl.handleWebError(w, http.StatusNotFound,
			fmt.Sprintf("There is no customer with id %d in the database", id))

		return
	}

	data, err := json.Marshal(customer)

	if err != nil {
		ctl.handleInternalError("Couldn't marshal data to JSON", err)
		ctl.handleWebError(w, http.StatusInternalServerError, "Couldn't marshal data to JSON")

		return
	}

	ctl.sendData(w, data)
}

func (ctl *CustomerController) getCustomers(w http.ResponseWriter, r *http.Request) {
	customers, err := ctl.customerRepo.GetAllCustomers()

	if err != nil {
		ctl.handleInternalError("Database access error", err)
		ctl.handleWebError(w, http.StatusNotFound,
			"Couldn't extract any entry from the customers database")

		return
	}

	data, err := json.Marshal(customers)

	if err != nil {
		ctl.handleInternalError("Couldn't marshal data to JSON", err)
		ctl.handleWebError(w, http.StatusInternalServerError, "Couldn't marshal data to JSON")

		return
	}

	ctl.sendData(w, data)
}

func (ctl *CustomerController) addCustomer(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)

	if err != nil {
		ctl.handleWebError(w, http.StatusBadRequest,
			"Couldn't read body")

		return
	}

	customer := new(repo.Customer)
	err = json.Unmarshal(data, customer)

	if err != nil {
		ctl.handleWebError(w, http.StatusBadRequest,
			"Couldn't parse JSON data")

		return
	}

	err = ctl.customerRepo.AddCustomer(customer)

	if err != nil {
		ctl.handleInternalError("Couldn't add customer to the database", err)
		ctl.handleWebError(w, http.StatusInternalServerError,
			"Couldn't add data to the database")

		return
	}

	ctl.sendSuccess(w, "Added successfully")
}

func (ctl *CustomerController) updateCustomer(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)

	if err != nil {
		ctl.handleWebError(w, http.StatusBadRequest,
			"Couldn't read body")

		return
	}

	customer := new(repo.Customer)
	err = json.Unmarshal(data, customer)

	if err != nil {
		ctl.handleWebError(w, http.StatusBadRequest,
			"Couldn't parse JSON data")

		return
	}

	// Check if exists.
	_, err = ctl.customerRepo.GetCustomerByID(customer.ID)

	if err != nil {
		ctl.handleInternalError("Database access error", err)
		ctl.handleWebError(w, http.StatusNotFound, "The customer doesn't exist")

		return
	}

	err = ctl.customerRepo.UpdateCustomer(customer)

	if err != nil {
		ctl.handleInternalError("Couldn't update customer in the database", err)
		ctl.handleWebError(w, http.StatusInternalServerError,
			"Couldn't update data in the database")

		return
	}

	ctl.sendSuccess(w, "Updated successfully")
}

func (ctl *CustomerController) deleteCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		ctl.handleWebError(w, http.StatusBadRequest,
			fmt.Sprintf("Incorrect parameter for id: %v", params["id"]))

		return
	}

	// Check if exists.
	_, err = ctl.customerRepo.GetCustomerByID(int64(id))

	if err != nil {
		ctl.handleInternalError("Database access error", err)
		ctl.handleWebError(w, http.StatusNotFound,
			"The customer doesn't exist")

		return
	}

	err = ctl.customerRepo.DeleteCustomer(int64(id))

	if err != nil {
		ctl.handleInternalError("Couldn't delete the customer", err)
		ctl.handleWebError(w, http.StatusInternalServerError,
			"Couldn't delete the customer from the database")

		return
	}

	ctl.sendSuccess(w, "Deleted successfully")
}

// SetupRoutes sets up routes for the controller.
func (ctl *CustomerController) SetupRoutes(router *mux.Router) {
	router.Use(jsonMiddleware)

	router.HandleFunc("/{id:[0-9]+}", ctl.getCustomer).Methods("GET")
	router.HandleFunc("/", ctl.getCustomers).Methods("GET")
	router.HandleFunc("/", ctl.addCustomer).Methods("POST")
	router.HandleFunc("/", ctl.updateCustomer).Methods("PATCH")
	router.HandleFunc("/{id:[0-9]+}", ctl.deleteCustomer).Methods("DELETE")
}

// NewCustomerController returns a new controller for the REST API operations on customers.
func NewCustomerController(repository repo.ICustomerRepository, logger *log.Logger) *CustomerController {
	ctl := new(CustomerController)

	ctl.customerRepo = repository
	ctl.logger = logger

	return ctl
}
