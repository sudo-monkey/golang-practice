package main

import (
	"fmt"
	"log"
	"net/http"

	handler "rest/pkg/handlers"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	//router
	r.HandleFunc("/", handler.SayHello("Shane")).Methods("GET")

	fmt.Println("server at port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))

}
