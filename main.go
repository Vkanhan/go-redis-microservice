package main

import (
	"context"
	"fmt"

	"github.com/Vkanhan/go-redis-microservice/application"
)

func main() {

	//create an instance of the application
	app := application.New()

	// Start the application by calling the Start method with a background context
	err := app.Start(context.TODO())
	if err != nil {
		fmt.Println("failed to start app: ", err)
	}
}
