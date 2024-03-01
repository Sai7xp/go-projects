/*
* Created on 01 March 2024
* @author Sai Sumanth
 */
package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

var logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)

// entry point
func main() {
	fmt.Println("wf-collect-client Microservice Started")

	router := mux.NewRouter()

	/// POST api for collecting details of build
	router.HandleFunc("/api/v1/collect", collectDetailsHandler).Methods("POST")

	/// wrong route handler
	router.NotFoundHandler = http.HandlerFunc(routeNotFoundHandler)

	fmt.Println("Server started at port 3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}

func collectDetailsHandler(res http.ResponseWriter, req *http.Request) {

}

// 404 route not found handler
func routeNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	logger.Println("Route Not Found")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"message":"OOPS! Wrong Route"}`))
}
