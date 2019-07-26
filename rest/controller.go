package rest

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Controller is a wrapper for controllers with web methods registered in the router.
type Controller interface {
	SetupRoutes(router *mux.Router)
}

type controller struct {
	logger *log.Logger
}

func (ctl *controller) handleInternalError(message string, err error) {
	if err != nil {
		ctl.logger.Printf("Error occured: %s, %s\n", message, err)
	}
}

func (ctl *controller) handleWebError(w http.ResponseWriter, statusCode int, message string) {
	mes := fmt.Sprintf("%d - %s", statusCode, message)

	ctl.logger.Println("Sent error message to the client:", mes)
	http.Error(w, mes, statusCode)
}

func (ctl *controller) sendSuccess(w http.ResponseWriter, message string) {
	http.Error(w, message, http.StatusOK)
}

func (ctl *controller) sendData(w http.ResponseWriter, data []byte) {
	ctl.logger.Println("Sending data to the client:", string(data))
	_, err := w.Write(data)
	ctl.handleInternalError("Couldn't write data to the HTTP network stream", err)
}
