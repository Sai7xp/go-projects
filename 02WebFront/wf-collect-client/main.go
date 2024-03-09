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
	"log"
	"net/http"
	"os"

	"github.com/IBM/sarama"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

var logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime) // logger
var validate *validator.Validate                               // validator

const (
	kafkaTopic = "webfront-kafka"
)

func main() {
	fmt.Println("wf-collect-client Microservice Started")
	validate = validator.New(validator.WithRequiredStructEnabled())

	/// Create New Route
	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()

	/// POST api for collecting details of build
	router.HandleFunc("/collect", collectDetailsHandler).Methods("POST")
	router.HandleFunc("/build/{build_id}", showBuildEventsHandler).Methods("GET")

	/// wrong route handler
	router.NotFoundHandler = http.HandlerFunc(routeNotFoundHandler)
	fmt.Println("Server started at port 3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}

// /
// 🛠️ Handler for /api/v1/collect API (only POST is allowed)
// /
func collectDetailsHandler(w http.ResponseWriter, req *http.Request) {
	projecDetails := new(models.WebfrontCollectDetails)
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
	/// generate random build id
	buildId := utils.GenerateBuildId()

	// send details along with buildId to another microservice `wf-code-builder` via kafka
	var builderReqBody = map[string]interface{}{
		"build_id":           buildId,
		"project_github_url": projecDetails.ProjectGithubUrl,
		"build_command":      projecDetails.BuildCommand,
		"build_out_dir":      projecDetails.BuildOutDir,
	}
	buildDetailsBytes, _ := json.Marshal(builderReqBody)

	// 🗣️ add new message to kafka
	kafkaErr := pushBuildDetailsToKafka(buildDetailsBytes)
	if kafkaErr != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":  false,
			"message":  kafkaErr.Error(),
			"build_id": nil,
		})
		return
	}

	/// return success message to user along with the build id
	logger.Printf("New Build request!! Assigned buildId %s", buildId)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"message":  "New build request added to the queue successfully!✅",
		"build_id": buildId,
	})
}

func createKafkaProducer(url []string) (sarama.SyncProducer, error) {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll
	kafkaConfig.Producer.Retry.Max = 5

	producer, err := sarama.NewSyncProducer(url, kafkaConfig)
	if err != nil {
		return nil, err
	}
	/// return producer and no error
	return producer, nil
}

// sends new message to kafka
func pushBuildDetailsToKafka(buildDetailsBytes []byte) error {

	url := []string{"localhost:29092"}
	producer, err := createKafkaProducer(url)
	if err != nil {
		return err
	}
	defer producer.Close()

	message := &sarama.ProducerMessage{
		Topic: kafkaTopic,
		Value: sarama.ByteEncoder(buildDetailsBytes),
	}
	_, _, err = producer.SendMessage(message)
	if err != nil {
		return err
	}
	log.Printf("Produced new message in topic %s", kafkaTopic)
	return nil

}

// 404 route not found handler
func routeNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	logger.Println("Route Not Found")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"message":"OOPS! Wrong Route"}`))
}

// Get Build events handler
func showBuildEventsHandler(w http.ResponseWriter, r *http.Request) {

}
