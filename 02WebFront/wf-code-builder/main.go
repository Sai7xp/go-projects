/*
* Created on 01 March 2024
* @author Sai Sumanth
 */
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"

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
	kafkaTopic       = "webfront-kafka" // topic name
	dockerVolumeName = "wf-storage"
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
			// 1. Build the image using docker build command
			// 2. Clone Repository and generate the build folder
			// 3. Deploy on nginx server
			cloneRepoAndGenerateBuild(buildRequestDetails)
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

// / initiate the repo cloning & build generation process
func cloneRepoAndGenerateBuild(buildDetails BuildRequestDetails) {
	// cmd := exec.Command("docker", "volume", "create", "wf-storage")
	logger.Println("Creating docker build....")
	cmd := exec.Command("docker", "build", "-t", "wf_build_react_app:latest", ".")
	cmd.Dir = "./wf_build_react_app"

	if _, err := cmd.CombinedOutput(); err != nil {
		fmt.Println("❌ Failed to build the image: ", err)
	}
	logger.Println("✅ Docker Image Created Successfully")
	// Build docker image success ✅

	/// Next step is to clone the repository and generate the build
	logger.Println("🤞 Trying to clone Repository and Generate Build....")
	buildIdArg := fmt.Sprintf("BUILD_ID=%s", buildDetails.BuildId)
	cmd = exec.Command("docker", "run", "-e", buildIdArg, "-v", "wf-storage:/wf/storage", "wf_build_react_app:latest", "-p", buildDetails.ProjectGithubUrl, "-b", buildDetails.BuildCommand, "-o", buildDetails.BuildOutDir)
	scriptFileOutput, err := cmd.CombinedOutput()
	if err != nil {
		logger.Println("❌ Something went wrong while cloning the repository: ", err)
		fmt.Println(scriptFileOutput)
		return
	}
	fmt.Println("✅ Cloned Repository from git url and generated build")

	/// Final Step: Deploy the Build
	ports := []int{5000, 8080, 6060, 7070}
	for _, port := range ports {
		if isPortAvailable(port) {
			logger.Println("Deployment Started....")
			pathArg := fmt.Sprintf("/var/lib/docker/volumes/wf-storage/_data/%s:/usr/share/nginx/html", buildDetails.BuildId)
			portArg := fmt.Sprintf("%d:80", port)
			serverNameArg := fmt.Sprintf("server-%s", buildDetails.BuildId)

			deployCmd := exec.Command("docker", "run", "--rm", "-d", "-p", portArg, "--name", serverNameArg, "-v", "wf-storage:/mnt", "-v", pathArg, "nginx")
			fmt.Println(deployCmd.Args)
			deployOutput, err := deployCmd.CombinedOutput()
			if err != nil {
				logger.Println("❌ Something went wrong while Deploying: ", err)
				fmt.Println(deployOutput)
				return
			}
			logger.Printf("🥳 Deployed Successfully at port %d. Click here to view the deployed version http://localhost:%d", port, port)
			break
		} else {
			logger.Println("🔌 ", port, "is already in use. Trying different port")
		}
	}
}

func isPortAvailable(port int) bool {
	// Attempt to listen on the specified port
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		// Port is already in use
		return false
	}
	// Close the listener
	listener.Close()
	// Port is available
	return true
}