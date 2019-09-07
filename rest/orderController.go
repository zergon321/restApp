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

// OrderController provides REST API methods for orders and their services.
type OrderController struct {
	orderRepo   repo.IOrderRepository
	serviceRepo repo.IServiceRepository
	controller
}

func (ctl *OrderController) getOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		ctl.handleWebError(w, http.StatusBadRequest,
			fmt.Sprintf("Incorrect parameter for id: %v", params["id"]))

		return
	}

	order, err := ctl.orderRepo.GetOrderByID(int64(id))

	if err != nil {
		ctl.handleInternalError("Database access error", err)
		ctl.handleWebError(w, http.StatusNotFound,
			fmt.Sprintf("There is no order with id %d in the database", id))

		return
	}

	data, err := json.Marshal(order)

	if err != nil {
		ctl.handleInternalError("Couldn't marshal data to JSON", err)
		ctl.handleWebError(w, http.StatusInternalServerError, "Couldn't marshal data to JSON")

		return
	}

	ctl.sendData(w, data)
}

func (ctl *OrderController) getOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := ctl.orderRepo.GetAllOrders()

	if err != nil {
		ctl.handleInternalError("Database access error", err)
		ctl.handleWebError(w, http.StatusNotFound,
			"Couldn't extract any entry from the orders database")

		return
	}

	data, err := json.Marshal(orders)

	if err != nil {
		ctl.handleInternalError("Couldn't marshal data to JSON", err)
		ctl.handleWebError(w, http.StatusInternalServerError, "Couldn't marshal data to JSON")

		return
	}

	ctl.sendData(w, data)
}

func (ctl *OrderController) addOrder(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)

	if err != nil {
		ctl.handleWebError(w, http.StatusBadRequest,
			"Couldn't read body")

		return
	}

	order := new(repo.Order)
	err = json.Unmarshal(data, order)

	if err != nil {
		ctl.handleWebError(w, http.StatusBadRequest,
			"Couldn't parse JSON data")

		return
	}

	err = ctl.orderRepo.AddOrder(order)

	if err != nil {
		ctl.handleInternalError("Couldn't add service to the database", err)
		ctl.handleWebError(w, http.StatusInternalServerError,
			"Couldn't add data to the database")

		return
	}

	ctl.sendSuccess(w, "Added successfully")
}

func (ctl *OrderController) updateOrder(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)

	if err != nil {
		ctl.handleWebError(w, http.StatusBadRequest,
			"Couldn't read body")

		return
	}

	order := new(repo.Order)
	err = json.Unmarshal(data, order)

	if err != nil {
		ctl.handleWebError(w, http.StatusBadRequest,
			"Couldn't parse JSON data")

		return
	}

	// Check if exists.
	_, err = ctl.orderRepo.GetOrderByID(order.ID)

	if err != nil {
		ctl.handleInternalError("Database access error", err)
		ctl.handleWebError(w, http.StatusNotFound, "The order doesn't exist")

		return
	}

	err = ctl.orderRepo.UpdateOrder(order)

	if err != nil {
		ctl.handleInternalError("Couldn't update service in the database", err)
		ctl.handleWebError(w, http.StatusInternalServerError,
			"Couldn't update data in the database")

		return
	}

	ctl.sendSuccess(w, "Updated successfully")
}

func (ctl *OrderController) deleteOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		ctl.handleWebError(w, http.StatusBadRequest,
			fmt.Sprintf("Incorrect parameter for id: %v", params["id"]))

		return
	}

	// Check if exists.
	_, err = ctl.orderRepo.GetOrderByID(int64(id))

	if err != nil {
		ctl.handleInternalError("Database access error", err)
		ctl.handleWebError(w, http.StatusNotFound,
			"The order doesn't exist")

		return
	}

	err = ctl.orderRepo.DeleteOrder(int64(id))

	if err != nil {
		ctl.handleInternalError("Couldn't delete the customer", err)
		ctl.handleWebError(w, http.StatusInternalServerError,
			"Couldn't delete the customer from the database")

		return
	}

	ctl.sendSuccess(w, "Deleted successfully")
}

func (ctl *OrderController) getOrderService(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	orderID, err := strconv.Atoi(params["orderId"])

	if err != nil {
		ctl.handleWebError(w, http.StatusBadRequest,
			fmt.Sprintf("Incorrect parameter for id: %v", params["orderId"]))

		return
	}

	serviceID, err := strconv.Atoi(params["serviceId"])

	if err != nil {
		ctl.handleWebError(w, http.StatusBadRequest,
			fmt.Sprintf("Incorrect parameter for id: %v", params["serviceId"]))

		return
	}

	service, err := ctl.orderRepo.GetOrderServiceByID(int64(orderID), int64(serviceID))

	if err != nil {
		ctl.handleInternalError("Database access error", err)
		ctl.handleWebError(w, http.StatusNotFound,
			fmt.Sprintf("There is no service with id %d for order with id %d in the database", serviceID, orderID))

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

func (ctl *OrderController) getOrderServices(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	orderID, err := strconv.Atoi(params["orderId"])

	if err != nil {
		ctl.handleWebError(w, http.StatusBadRequest,
			fmt.Sprintf("Incorrect parameter for id: %v", params["orderId"]))

		return
	}

	services, err := ctl.orderRepo.GetAllOrderServices(int64(orderID))

	if err != nil {
		ctl.handleInternalError("Database access error", err)
		ctl.handleWebError(w, http.StatusNotFound,
			fmt.Sprintf("There are no services for order with id %d in the database", orderID))

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

func (ctl *OrderController) addOrderService(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	orderID, err := strconv.Atoi(params["orderId"])

	if err != nil {
		ctl.handleWebError(w, http.StatusBadRequest,
			fmt.Sprintf("Incorrect parameter for id: %v", params["orderId"]))

		return
	}

	serviceID, err := strconv.Atoi(params["serviceId"])

	if err != nil {
		ctl.handleWebError(w, http.StatusBadRequest,
			fmt.Sprintf("Incorrect parameter for id: %v", params["serviceId"]))

		return
	}

	// Check if the order exists.
	_, err = ctl.orderRepo.GetOrderByID(int64(orderID))

	if err != nil {
		ctl.handleInternalError("Database access error", err)
		ctl.handleWebError(w, http.StatusNotFound,
			"The order doesn't exist")

		return
	}

	// Check if the service exists.
	_, err = ctl.serviceRepo.GetServiceByID(int64(serviceID))

	if err != nil {
		ctl.handleInternalError("Database access error", err)
		ctl.handleWebError(w, http.StatusNotFound,
			"The service doesn't exist")

		return
	}

	err = ctl.orderRepo.AddServiceToOrder(int64(orderID), int64(serviceID))

	if err != nil {
		ctl.handleInternalError("Couldn't add service to the database", err)
		ctl.handleWebError(w, http.StatusInternalServerError,
			"Couldn't add data to the database")

		return
	}

	ctl.sendSuccess(w, "Added successfully")
}

func (ctl *OrderController) deleteOrderSevice(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	orderID, err := strconv.Atoi(params["orderId"])

	if err != nil {
		ctl.handleWebError(w, http.StatusBadRequest,
			fmt.Sprintf("Incorrect parameter for id: %v", params["orderId"]))

		return
	}

	serviceID, err := strconv.Atoi(params["serviceId"])

	if err != nil {
		ctl.handleWebError(w, http.StatusBadRequest,
			fmt.Sprintf("Incorrect parameter for id: %v", params["serviceId"]))

		return
	}

	// Check if the order exists.
	_, err = ctl.orderRepo.GetOrderByID(int64(orderID))

	if err != nil {
		ctl.handleInternalError("Database access error", err)
		ctl.handleWebError(w, http.StatusNotFound,
			"The order doesn't exist")

		return
	}

	// Check if the service exists in the order.
	_, err = ctl.orderRepo.GetOrderServiceByID(int64(orderID), int64(serviceID))

	if err != nil {
		ctl.handleInternalError("Database access error", err)
		ctl.handleWebError(w, http.StatusNotFound,
			"The service doesn't exist or isn't included in the order")

		return
	}

	err = ctl.orderRepo.DeleteServiceFromOrder(int64(orderID), int64(serviceID))

	if err != nil {
		ctl.handleInternalError("Couldn't delete the service from the order", err)
		ctl.handleWebError(w, http.StatusInternalServerError,
			"Couldn't delete the service from the order")

		return
	}

	ctl.sendSuccess(w, "Deleted successfully")
}

// SetupRoutes sets up routes for the controller.
func (ctl *OrderController) SetupRoutes(router *mux.Router) {
	router.Use(jsonMiddleware)

	router.HandleFunc("/{id:[0-9]+}", ctl.getOrder).Methods("GET")
	router.HandleFunc("/", ctl.getOrders).Methods("GET")
	router.HandleFunc("/", ctl.addOrder).Methods("POST")
	router.HandleFunc("/", ctl.updateOrder).Methods("PATCH")
	router.HandleFunc("/{id:[0-9]+}", ctl.deleteOrder).Methods("DELETE")

	router.HandleFunc("/{orderId:[0-9]+}/services/{serviceId:[0-9]+}",
		ctl.getOrderService).Methods("GET")
	router.HandleFunc("/{orderId:[0-9]}/services",
		ctl.getOrderServices).Methods("GET")
	router.HandleFunc("/{orderId:[0-9]+}/services/{serviceId:[0-9]+}",
		ctl.addOrderService).Methods("POST")
	router.HandleFunc("/{orderId:[0-9]+}/services/{serviceId:[0-9]+}",
		ctl.deleteOrderSevice).Methods("DELETE")
}

// NewOrderController returns a new controller for the REST API operations on orders.
func NewOrderController(orderRepository repo.IOrderRepository,
	serviceRepository repo.IServiceRepository, logger *log.Logger) *OrderController {
	ctl := new(OrderController)

	ctl.orderRepo = orderRepository
	ctl.serviceRepo = serviceRepository
	ctl.logger = logger

	return ctl
}
