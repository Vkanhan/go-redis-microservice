package application

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

//App struct stores the router to the http.Handler
type App struct {
	router http.Handler	// router for handling HTTP requests
	rdb *redis.Client	// Redis client for database operations
}

// New creates a new instance of the App struct, initializes the router, and returns the instance.
func New() *App {
	// Initialize the App struct with the router configured by loadRoutes
	app := &App{
		router: loadRoutes(),					  // Initialize router
		rdb: redis.NewClient(&redis.Options{}),	 // Initialize Redis client
	}
	return app
}

// Start starts the HTTP server with the given context.
func (a *App) Start(ctx context.Context) error {
	server := http.Server{
		Addr: ":3000",
		Handler: a.router,
	}

	// Ping the Redis server to ensure connection
	err := a.rdb.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	// Ensure the Redis client is closed when the function exits
	defer func() {
		if err := a.rdb.Close(); err != nil {
			fmt.Println("failed to close redis", err)
		}
	}()

	fmt.Println("Starting server")

	// Channel to capture server errors
	ch := make(chan error, 1)

	go func() {
		// Start the HTTP server
		err = server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to start server: %w", err)
		}
		close(ch)

	}()

	// Select statement to handle either server error or context cancellation
	select {
	case err = <- ch:
		return err
	case <- ctx.Done():
		// If the context is cancelled, initiate server shutdown
		// Create a new context with a timeout for the shutdown process
		timeout, cancel := context.WithTimeout(context.Background(), time.Second * 10)
		defer cancel()

		//shut down the server gracefully
		return server.Shutdown(timeout)
	}
	
}

