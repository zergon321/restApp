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

// ServiceController provides REST API methods for services.
type ServiceController struct {
	serviceRepo repo.IServiceRepository
	controller
}

func (ctl *ServiceController) getService(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		ctl.handleWebError(w, http.StatusBadRequest,
			fmt.Sprintf("Incorrect parameter for id: %v", params["id"]))

		return
	}

	service, err := ctl.serviceRepo.GetServiceByID(int64(id))

	if err != nil {
		ctl.handleInternalError("Database access error", err)
		ctl.handleWebError(w, http.StatusNotFound,
			fmt.Sprintf("There is no service with id %d in the database", id))

		return
	}

	data, err := json.Marshal(service)

	if err != nil {
		ctl.handleInternalError("Couldn't marshal data to JSON", err)
		ctl.handleWebError(w, http.StatusInternalServerError, "Couldn't marshal data to JSON")

		return
	}

	ctl.sendData(w, data)
}

func (ctl *ServiceController) getServices(w http.ResponseWriter, r *http.Request) {
	services, err := ctl.serviceRepo.GetAllServices()

	if err != nil {
		ctl.handleInternalError("Database access error", err)
		ctl.handleWebError(w, http.StatusNotFound,
			"Couldn't extract any entry from the services database")

		return
	}

	data, err := json.Marshal(services)

	if err != nil {
		ctl.handleInternalError("Couldn't marshal data to JSON", err)
		ctl.handleWebError(w, http.StatusInternalServerError, "Couldn't marshal data to JSON")

		return
	}

	ctl.sendData(w, data)
}

func (ctl *ServiceController) addService(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)

	if err != nil {
		ctl.handleWebError(w, http.StatusBadRequest,
			"Couldn't read body")

		return
	}

	service := new(repo.Service)
	err = json.Unmarshal(data, service)

	if err != nil {
		ctl.handleWebError(w, http.StatusBadRequest,
			"Couldn't parse JSON data")

		return
	}

	err = ctl.serviceRepo.AddService(service)

	if err != nil {
		ctl.handleInternalError("Couldn't add service to the database", err)
		ctl.handleWebError(w, http.StatusInternalServerError,
			"Couldn't add data to the database")

		return
	}

	ctl.sendSuccess(w, "Added successfully")
}

func (ctl *ServiceController) updateService(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)

	if err != nil {
		ctl.handleWebError(w, http.StatusBadRequest,
			"Couldn't read body")

		return
	}

	service := new(repo.Service)
	err = json.Unmarshal(data, service)

	if err != nil {
		ctl.handleWebError(w, http.StatusBadRequest,
			"Couldn't parse JSON data")

		return
	}

	// Check if exists.
	_, err = ctl.serviceRepo.GetServiceByID(service.ID)

	if err != nil {
		ctl.handleInternalError("Database access error", err)
		ctl.handleWebError(w, http.StatusNotFound, "The service doesn't exist")

		return
	}

	err = ctl.serviceRepo.UpdateService(service)

	if err != nil {
		ctl.handleInternalError("Couldn't update service in the database", err)
		ctl.handleWebError(w, http.StatusInternalServerError,
			"Couldn't update data in the database")

		return
	}

	ctl.sendSuccess(w, "Updated successfully")
}

func (ctl *ServiceController) deleteService(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		ctl.handleWebError(w, http.StatusBadRequest,
			fmt.Sprintf("Incorrect parameter for id: %v", params["id"]))

		return
	}

	// Check if exists.
	_, err = ctl.serviceRepo.GetServiceByID(int64(id))

	if err != nil {
		ctl.handleInternalError("Database access error", err)
		ctl.handleWebError(w, http.StatusNotFound,
			"The service doesn't exist")

		return
	}

	err = ctl.serviceRepo.DeleteService(int64(id))

	if err != nil {
		ctl.handleInternalError("Couldn't delete the service", err)
		ctl.handleWebError(w, http.StatusInternalServerError,
			"Couldn't delete the service from the database")

		return
	}

	ctl.sendSuccess(w, "Deleted successfully")
}

// SetupRoutes sets up routes for the controller.
func (ctl *ServiceController) SetupRoutes(router *mux.Router) {
	router.Use(jsonMiddleware)

	router.HandleFunc("/{id:[0-9]+}", ctl.getService).Methods("GET")
	router.HandleFunc("/", ctl.getServices).Methods("GET")
	router.HandleFunc("/", ctl.addService).Methods("POST")
	router.HandleFunc("/", ctl.updateService).Methods("PATCH")
	router.HandleFunc("/services/{id:[0-9]+}", ctl.deleteService).Methods("DELETE")
}

// NewServiceController returns a new controller for the REST API operations on services.
func NewServiceController(repository repo.IServiceRepository, logger *log.Logger) *ServiceController {
	ctl := new(ServiceController)

	ctl.serviceRepo = repository
	ctl.logger = logger

	return ctl
}
