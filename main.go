package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"restApp/conf"
	"restApp/repo"
	"restApp/rest"

	"github.com/gorilla/mux"

	"gopkg.in/yaml.v2"

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

	// Create log file and logger.
	file, err := os.OpenFile(LOG, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalln("Couldn't open log file:", err)
	}
	defer file.Close()

	stream := io.MultiWriter(os.Stdout, file)
	logger := log.New(stream, PREFIX, log.LstdFlags|log.Lshortfile)

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
	db, err := sql.Open(dbDriver, fmt.Sprintf("%s://%s:%s@%s/%s",
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
