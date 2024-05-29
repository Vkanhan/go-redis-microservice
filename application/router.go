package application

import (
	"net/http"

	"github.com/Vkanhan/go-redis-microservice/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// loadRoutes sets up the routes for the application and returns the configured router.
func loadRoutes() *chi.Mux {
	// Create a new router instance using chi
	router := chi.NewRouter()

	// Use the Logger middleware to log
	router.Use(middleware.Logger)

	// Define a route for the root path "/" that responds with HTTP 200 status
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	//group all the routes under /order handler
	router.Route("/orders", loadOrderRoutes)
	
	// Return the configured router
	return router
}

// loadOrderRoutes sets up the routes related to order operations
func loadOrderRoutes(router chi.Router) {

	// Initialize an order handler
	orderHandler := &handler.Order{}

	router.Post("/", orderHandler.Create)
	router.Get("/", orderHandler.List)
	router.Get("/{id}", orderHandler.GetByID)
	router.Put("/{id}", orderHandler.UpdateByID)
	router.Delete("/{id}", orderHandler.DeleteByID)

}