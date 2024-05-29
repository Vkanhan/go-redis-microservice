package handler

import (
	"fmt"
	"net/http"
)

// Order struct is a handler for order-related HTTP requests
type Order struct {

}

// Create handles the HTTP POST request to create a new order
func (o *Order) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create an order")
}

// List handles the HTTP GET request to list all orders
func (o *Order) List(w http.ResponseWriter, r *http.Request) {
	fmt.Println("List all orders")
}

// GetByID handles the HTTP GET request to retrieve an order by its ID
func (o *Order) GetByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get the order by ID")
}

// UpdateByID handles the HTTP PUT request to update an order by its ID
func (o *Order) UpdateByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update the order by ID")
}

// DeleteByID handles the HTTP DELETE request to delete an order by its ID
func (o *Order) DeleteByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete the order by ID")
}
