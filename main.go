package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"restApp/conf"
	"restApp/repo"
	"restApp/rest"

	"gopkg.in/yaml.v2"

	_ "github.com/lib/pq"
)

// Configuration constants for the application.
const (
	LOG    = "sys.log"
	CONFIG = "config.yml"
	PREFIX = "app: "
)

func main() {
	err := os.Chdir("bin")

	if err != nil {
		log.Fatalln("Couldn't change directory to bin")
	}

	// Create log file and logger.
	file, err := os.OpenFile(LOG, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalln("Couldn't open log file:", err)
	}
	defer file.Close()

	logger := log.New(file, PREFIX, log.LstdFlags|log.Lshortfile)

	// Retrieve configuration settings from the yml-file.
	confFile, err := os.OpenFile(CONFIG, os.O_RDONLY, 0666)

	if err != nil {
		logger.Fatalln("Couldn't open config file:", err)
	}

	data, err := ioutil.ReadAll(confFile)

	if err != nil {
		logger.Fatalln("Couldn't read config file:", err)
	}

	config := conf.DbConfiguration{}
	err = yaml.Unmarshal(data, &config)

	if err != nil {
		logger.Fatalln("Couldn't parse config file:", err)
	}

	// Open database connection.
	db, err := sql.Open(config.Driver, fmt.Sprintf("%s://%s:%s@%s/%s",
		config.Protocol, config.Username, config.Password, config.Host, config.DbName))

	if err != nil {
		logger.Fatalln("Couldn't establish a db connection:", err)
	}
	defer db.Close()

	// Create data repositories.
	customerRepo := repo.NewCustomerRepo(db)
	serviceRepo := repo.NewServiceRepo(db)
	orderRepo := repo.NewOrderRepository(db)

	// Create REST API controllers.
	customerController := rest.NewCustomerController(customerRepo, logger)
	serviceController := rest.NewServiceController(serviceRepo, logger)
	orderController := rest.NewOrderController(orderRepo, serviceRepo, logger)

	// Get control muxes.
	mainMux := http.NewServeMux()
	customerMux := customerController.GetRouter()
	serviceMux := serviceController.GetRouter()
	orderMux := orderController.GetRouter()

	mainMux.Handle("/customers/", customerMux)
	mainMux.Handle("/services/", serviceMux)
	mainMux.Handle("/orders/", orderMux)

	http.ListenAndServe(":80", mainMux)
}
