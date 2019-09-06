package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"restApp/repo"
	"restApp/rest"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

// Configuration constants for the application.
const (
	LOG    = "sys.log"
	CONFIG = "config.yml"
	PREFIX = "app: "
)

var (
	dbDriver   string
	dbProtocol string
	dbUsername string
	dbPassword string
	dbHost     string
	dbName     string

	address string
	port    string
)

// parseFlags parses command line arguments and assigns them to global variables.
func parseFlags() {
	flag.StringVar(&dbDriver, "dbdriver", "", "A driver to access the database")
	flag.StringVar(&dbProtocol, "dbprotocol", "", "A protocol to access the database")
	flag.StringVar(&dbUsername, "dbusername", "", "A username to access the database")
	flag.StringVar(&dbPassword, "dbpassword", "", "A password to access the database")
	flag.StringVar(&dbHost, "dbhost", "", "A host on which the DBMS is deployed")
	flag.StringVar(&dbName, "dbname", "", "A name of the database")

	flag.StringVar(&address, "address", "", "An address to listen on")
	flag.StringVar(&port, "port", "80", "A port to listen on")

	flag.Parse()
}

func main() {
	parseFlags()

	// Change working directory to the application directory.
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))

	if err != nil {
		log.Fatalln("Couldn't get current application folder path:", err)
	}

	err = os.Chdir(dir)

	if err != nil {
		log.Fatalln("Couldn't change directory to bin:", err)
	}

	// Create a log file and a logger.
	file, err := os.OpenFile(LOG, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalln("Couldn't open log file:", err)
	}
	defer file.Close()

	stream := io.MultiWriter(os.Stdout, file)
	logger := log.New(stream, PREFIX, log.LstdFlags|log.Lshortfile)

	// Open database connection.
	db, err := sql.Open(dbDriver, fmt.Sprintf("%s://%s:%s@%s/%s?sslmode=disable",
		dbProtocol, dbUsername, dbPassword, dbHost, dbName))

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

	// Setup REST routes.
	router := mux.NewRouter()
	customers := router.PathPrefix("/customers").Subrouter()
	services := router.PathPrefix("/services").Subrouter()
	orders := router.PathPrefix("/orders").Subrouter()

	customerController.SetupRoutes(customers)
	serviceController.SetupRoutes(services)
	orderController.SetupRoutes(orders)

	addr := fmt.Sprintf("%s:%s", address, port)
	http.ListenAndServe(addr, router)
}
