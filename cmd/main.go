package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/balgabekj/go-ecommerce/pkg/model"
	"github.com/gorilla/mux"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type config struct {
	port string
	env  string
	db   struct {
		dsn string
	}
}

type application struct {
	config config
	models model.Models
}

func main() {
	var cfg config
	flag.StringVar(&cfg.port, "port", ":8081", "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:postgres@localhost/go_ecommerce?sslmode=disable", "PostgreSQL DSN")
	flag.Parse()
	db, err := openDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := &application{
		config: cfg,
		models: model.NewModels(db),
	}
	app.run()
}
func (app *application) run() {
	r := mux.NewRouter()

	v1 := r.PathPrefix("/api/v1").Subrouter()

	v1.HandleFunc("/api/v1/users", app.getAllUsersHandler).Methods("GET")            // Get all Users
	v1.HandleFunc("/api/v1/users/{userID}", app.getUserByIdHandler).Methods("GET")   // Get a specific product
	v1.HandleFunc("/api/v1/users", app.createUserHandler).Methods("POST")            // Create a new product
	v1.HandleFunc("/api/v1/users/{userID}", app.updateUserHandler).Methods("PUT")    // Update a product
	v1.HandleFunc("/api/v1/users/{userID}", app.deleteUserHandler).Methods("DELETE") // Delete a product

	// Product routes
	//v1.HandleFunc("/api/v1/products", app.listProductsHandler).Methods("GET")                 // Get all products
	v1.HandleFunc("/api/v1/products/{productID}", app.getProductHandler).Methods("GET")       // Get a specific product
	v1.HandleFunc("/api/v1/products", app.createProductHandler).Methods("POST")               // Create a new product
	v1.HandleFunc("/api/v1/products/{productID}", app.updateProductHandler).Methods("PUT")    // Update a product
	v1.HandleFunc("/api/v1/products/{productID}", app.deleteProductHandler).Methods("DELETE") // Delete a product

	// Order routes
	//v1.HandleFunc("/api/v1/orders", app.listOrdersHandler).Methods("GET")               // Get all orders
	v1.HandleFunc("/api/v1/orders/{orderID}", app.getOrderHandler).Methods("GET")       // Get a specific order
	v1.HandleFunc("/api/v1/orders", app.createOrderHandler).Methods("POST")             // Create a new order
	v1.HandleFunc("/api/v1/orders/{orderID}", app.updateOrderHandler).Methods("PUT")    // Update an order
	v1.HandleFunc("/api/v1/orders/{orderID}", app.deleteOrderHandler).Methods("DELETE") // Delete an order

	log.Printf("Starting server on %s\n", app.config.port)
	err := http.ListenAndServe(app.config.port, r)
	log.Fatal(err)
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %v", err)
	}

	log.Println("Successfully connected to the database")
	return db, nil
}
