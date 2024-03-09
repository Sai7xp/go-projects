/*
* Created on 01 March 2024
* @author Sai Sumanth
 */
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/IBM/sarama"
)

var logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime) // logger

// build request received via kafka
type BuildRequestDetails struct {
	BuildId          string `json:"build_id"`
	ProjectGithubUrl string `json:"project_github_url"`
	BuildCommand     string `json:"build_command"`
	BuildOutDir      string `json:"build_out_dir"`
}

const (
	kafkaTopic = "webfront-kafka" // topic name
)

func main() {
	fmt.Println("wf-code-builder Microservice Started")

	/// start listening to kafka messages produced on 'webfront-kafka' topic
	if con, err := createNewConsumer(); err == nil {
		defer con.Close()

		partitionConsumer, err := con.ConsumePartition(kafkaTopic, 0, sarama.OffsetNewest)
		if err != nil {
			log.Fatalf("Error creating partition consumer: %v", err)
		}
		defer partitionConsumer.Close()

		for msg := range partitionConsumer.Messages() {
			var buildRequestDetails BuildRequestDetails
			err := json.Unmarshal(msg.Value, &buildRequestDetails)
			if err != nil {
				log.Printf("Error while parsing message: %v", err)
				continue
			}

			logger.Printf("📣 Received New Message in topic %s: %+v\n", msg.Topic, buildRequestDetails)
		}
	}
}

func createNewConsumer() (sarama.Consumer, error) {
	config := sarama.NewConfig()
	consumer, err := sarama.NewConsumer([]string{"localhost:29092"}, config)
	if err != nil {
		log.Fatalf("Error creating Kafka consumer: %v", err)
		return nil, err
	}
	return consumer, nil

}
