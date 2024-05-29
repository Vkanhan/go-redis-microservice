package application

import (
	"context"
	"fmt"
	"net/http"
)

//App struct stores the router to the http.Handler
type App struct {
	router http.Handler
}

// New creates a new instance of the App struct, initializes the router, and returns the instance.
func New() *App {
	// Initialize the App struct with the router configured by loadRoutes
	app := &App{
		router: loadRoutes(),
	}
	return app
}

// Start starts the HTTP server with the given context.
func (a *App) Start(ctx context.Context) error {
	server := http.Server{
		Addr: ":3000",
		Handler: a.router,
	}

	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}

