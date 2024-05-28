package main 

import(
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/middleware"

)

func main() {
	//create a new router using chi library
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	//define a route that match the path "/hello" and binds it to basicHandler func
	router.Get("/hello", basicHandler)

	//create an HTTP server instance with specified address and router
	server := &http.Server{
		Addr: ":3000",
		Handler: router,
	}

	//start the server and listen to incoming request
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("failed to listen to port :3000", err)
	}
}

//basicHandler  func will be called when request is made to "/hello" route
func basicHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("go-microservice"))
}