package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/Vkanhan/go-redis-microservice/application"
)

func main() {

	//create an instance of the application
	app := application.New(application.LoadConfig())

	// Create a context that listens for an interrupt signal and cancel function is called when main exist
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Start the application by calling the Start method with a background context
	err := app.Start(ctx)
	if err != nil {
		fmt.Println("failed to start app: ", err)
	}
}
