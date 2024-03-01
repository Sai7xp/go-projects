/*
* Created on 01 March 2024
* @author Sai Sumanth
 */

package main

import (
	"collectclient/models"
	"collectclient/utils"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

var logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)

// validator
var validate *validator.Validate

// entry point
func main() {
	fmt.Println("wf-collect-client Microservice Started")

	router := mux.NewRouter()
	validate = validator.New(validator.WithRequiredStructEnabled())

	/// POST api for collecting details of build
	router.HandleFunc("/api/v1/collect", collectDetailsHandler).Methods("POST")
	/// wrong route handler
	router.NotFoundHandler = http.HandlerFunc(routeNotFoundHandler)

	fmt.Println("Server started at port 3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}

// 🛠️ Handler for /api/v1/collect API (allowed methods POST)
func collectDetailsHandler(w http.ResponseWriter, req *http.Request) {
	var projecDetails models.WebfrontCollectDetails
	/// decode data
	if err := json.NewDecoder(req.Body).Decode(&projecDetails); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	/// validation
	if err := validate.Struct(projecDetails); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	buildId := utils.GenerateBuildId()

	var response = make(map[string]interface{})
	response["success"] = true
	response["buildId"] = buildId
	logger.Printf("New Build request!! Assigned buildId %s", buildId)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	// TODO: send details along with buildId to another microservice `wf-code-builder`
	// sendDetailsToCodeBuilder()
}

// 404 route not found handler
func routeNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	logger.Println("Route Not Found")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"message":"OOPS! Wrong Route"}`))
}
