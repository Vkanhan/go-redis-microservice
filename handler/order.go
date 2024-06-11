package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/Vkanhan/go-redis-microservice/model"
	"github.com/Vkanhan/go-redis-microservice/repository/order"
)

// Order struct is a handler for order-related HTTP requests
type Order struct {
	Repo *order.RedisRepo
}

// Create handles the HTTP POST request to create a new order
func (h *Order) Create(w http.ResponseWriter, r *http.Request) {
	// Decode the request body into a struct
	var body struct {
		CustomerID uuid.UUID         `json:"customer_id"`
		LineItems  []model.LineItems `json:"line_items"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the current UTC time
	now := time.Now().UTC()

	// Create a new order instance
	order := model.Order{
		OrderID:    uint64(uuid.New().ID()),
		CustomerID: body.CustomerID,
		LineItems:  body.LineItems,
		CreatedAt:  &now,
	}

	// Insert the new order into the repository
	err := h.Repo.Insert(r.Context(), order)
	if err != nil {
		fmt.Println("failed to insert: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Set the content type header
	w.Header().Set("Content-Type", "application/json")

	// Set the status code to 201 Created
	w.WriteHeader(http.StatusCreated)

	//Encode the order response
	if err := json.NewEncoder(w).Encode(order); err != nil {
		fmt.Println("failed to encode the order: ", err)
		return
	}

}

// List handles the HTTP GET request to list all orders
func (h *Order) List(w http.ResponseWriter, r *http.Request) {
	// Get the cursor parameter from the URL query, default to 0
	cursorStr := r.URL.Query().Get("cursor")
	if cursorStr == "" {
		cursorStr = "0"
	}

	const decimal = 10
	const bitSize = 64

	// Parse the cursor string to an unsigned integer
	cursor, err := strconv.ParseUint(cursorStr, decimal, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Define the size of the result set
	const size = 50

	// Find all orders from the repository starting from the cursor
	res, err := h.Repo.FindAll(r.Context(), order.FindAllPage{
		Offset: cursor,
		Size:   size,
	})
	if err != nil {
		fmt.Println("failed to find all:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Prepare the response structure
	var response struct {
		Items []model.Order `json:"items"`
		Next  uint64        `json:"next,omitempty"`
	}
	response.Items = res.Orders
	response.Next = res.Cursor

	// Set the content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Set the status code to 201 Created
	w.WriteHeader(http.StatusOK)

	// Encode the response into JSON and write it to the response writer
	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Println("failed to encode response:", err)
		return
	}

}

// GetByID handles the HTTP GET request to retrieve an order by its ID
func (h *Order) GetByID(w http.ResponseWriter, r *http.Request) {
	// Get the order ID from the URL parameter and parse it
	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64

	orderID, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Find the order by ID from the repository
	o, err := h.Repo.FindByID(r.Context(), orderID)
	if errors.Is(err, order.ErrorNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("failed to find by id:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Encode the order to JSON and write the response
	if err := json.NewEncoder(w).Encode(o); err != nil {
		fmt.Println("failed to marshal:", err)
		return
	}

}

// UpdateByID handles the HTTP PUT request to update an order by its ID
func (h *Order) UpdateByID(w http.ResponseWriter, r *http.Request) {
	// Decode the request body into a struct
	var body struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the order ID from the URL parameter and parse it
	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64

	orderID, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Find the order by ID from the repository
	theOrder, err := h.Repo.FindByID(r.Context(), orderID)
	if errors.Is(err, order.ErrorNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("failed to find by id:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	const completedStatus = "completed"
	const shippedStatus = "shipped"

	// Update the order status based on the request body
	now := time.Now().UTC()

	switch body.Status {
	case shippedStatus:
		if theOrder.ShippedAT != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		theOrder.ShippedAT = &now
	case completedStatus:
		if theOrder.CompletedAt != nil || theOrder.ShippedAT == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		theOrder.CompletedAt = &now
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Update the order in the repository
	err = h.Repo.UpdateByID(r.Context(), theOrder)
	if err != nil {
		fmt.Println("failed to insert:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Encode the updated order to JSON and write the response
	if err := json.NewEncoder(w).Encode(theOrder); err != nil {
		fmt.Println("failed to marshal:", err)
		return
	}
}

// DeleteByID handles the HTTP DELETE request to delete an order by its ID
func (h *Order) DeleteByID(w http.ResponseWriter, r *http.Request) {
	// Get the order ID from the URL parameter and parse it
	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64

	orderID, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Delete the order by ID from the repository
	err = h.Repo.DeleteByID(r.Context(), orderID)
	if errors.Is(err, order.ErrorNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("failed to find by id:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
