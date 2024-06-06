package application

import (
	"net/http"

	"github.com/Vkanhan/go-redis-microservice/handler"
	"github.com/Vkanhan/go-redis-microservice/repository/order"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// loadRoutes sets up the routes for the application and returns the configured router.
func (a *App) loadRoutes() {
	// Create a new router instance using chi
	router := chi.NewRouter()

	// Use the Logger middleware to log
	router.Use(middleware.Logger)

	// Define a route for the root path "/" that responds with HTTP 200 status
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	//group all the routes under /order handler
	router.Route("/orders", a.loadOrderRoutes)

	// Return the configured router
	a.router = router
}

// loadOrderRoutes sets up the routes related to order operations
func (a *App)loadOrderRoutes(router chi.Router) {

	// Initialize an order handler
	orderHandler := &handler.Order{
		Repo: &order.RedisRepo{
			Client: a.rdb,
		},
	}

	router.Post("/", orderHandler.Create)
	router.Get("/", orderHandler.List)
	router.Get("/{id}", orderHandler.GetByID)
	router.Put("/{id}", orderHandler.UpdateByID)
	router.Delete("/{id}", orderHandler.DeleteByID)

}
