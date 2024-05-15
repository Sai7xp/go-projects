/*
* Created on 01 March 2024
* @author Sai Sumanth
 */
package main

import (
	"code-builder-service/database"
	"code-builder-service/models"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"time"

	"github.com/IBM/sarama"
)

var logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime) // logger

const (
	kafkaTopic       = "webfront-kafka" // topic name
	dockerVolumeName = "wf-storage"
)

func main() {
	fmt.Println("wf-code-builder Microservice Started")
	database.Init()

	// Build Docker Image on Init
	logger.Println("Creating docker build....")
	cmd := exec.Command("docker", "build", "-t", "wf_build_react_app:latest", ".")
	cmd.Dir = "./wf_build_react_app"

	if _, err := cmd.CombinedOutput(); err != nil {
		log.Fatal("❌ Failed to build the image: ", err)
	}
	logger.Println("✅ Docker Image Created Successfully")

	/// start listening to kafka messages produced on 'webfront-kafka' topic
	if con, err := createNewConsumer(); err == nil {
		defer con.Close()

		partitionConsumer, err := con.ConsumePartition(kafkaTopic, 0, sarama.OffsetNewest)
		if err != nil {
			log.Fatalf("Error creating partition consumer: %v", err)
		}
		defer partitionConsumer.Close()

		for msg := range partitionConsumer.Messages() {
			var buildRequestDetails models.BuildRequestDetails
			err := json.Unmarshal(msg.Value, &buildRequestDetails)
			if err != nil {
				log.Printf("Error while parsing message: %v", err)
				continue
			}

			logger.Printf("📣 Received New Message in topic %s: %+v\n", msg.Topic, buildRequestDetails)
			database.AddNewBuild(buildRequestDetails)
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
func cloneRepoAndGenerateBuild(buildDetails models.BuildRequestDetails) error {
	// cmd := exec.Command("docker", "volume", "create", "wf-storage")
	
	// npm build started
	database.UpdateEventToExistingBuild(buildDetails.BuildId, map[string]map[string]interface{}{"BUILD_STARTED": {"timestamp": time.Now()}})

	/// Next step is to clone the repository and generate the build
	logger.Println("🤞 Trying to clone Repository and Generate Build....")
	buildIdArg := fmt.Sprintf("BUILD_ID=%s", buildDetails.BuildId)
	cmd := exec.Command("docker", "run", "-e", buildIdArg, "-v", "wf-storage:/wf/storage", "wf_build_react_app:latest", "-p", buildDetails.ProjectGithubUrl, "-b", buildDetails.BuildCommand, "-o", buildDetails.BuildOutDir)
	scriptFileOutput, err := cmd.CombinedOutput()
	if err != nil {
		logger.Println("❌ Something went wrong while cloning the repository: ", err)
		fmt.Println(string(scriptFileOutput))
		database.UpdateEventToExistingBuild(buildDetails.BuildId, map[string]map[string]interface{}{"BUILD_FAILED": {"timestamp": time.Now(), "reason": string(scriptFileOutput)}})
		return err
	}
	fmt.Println("✅ Cloned Repository from git url and generated build")
	database.UpdateEventToExistingBuild(buildDetails.BuildId, map[string]map[string]interface{}{"BUILD_PASSED": {"timestamp": time.Now()}})

	/// Final Step: Deploy the Build
	ports := []int{5000, 8080, 6060, 7070}
	for _, port := range ports {
		if isPortAvailable(port) {
			logger.Println("Deployment Started....")
			database.UpdateEventToExistingBuild(buildDetails.BuildId, map[string]map[string]interface{}{"DEPLOY_STARTED": {"timestamp": time.Now()}})

			pathArg := fmt.Sprintf("/var/lib/docker/volumes/wf-storage/_data/%s:/usr/share/nginx/html", buildDetails.BuildId)
			portArg := fmt.Sprintf("%d:80", port)
			serverNameArg := fmt.Sprintf("server-%s", buildDetails.BuildId)

			deployCmd := exec.Command("docker", "run", "--rm", "-d", "-p", portArg, "--name", serverNameArg, "-v", "wf-storage:/mnt", "-v", pathArg, "nginx")
			fmt.Println(deployCmd.Args)
			deployOutput, err := deployCmd.CombinedOutput()
			if err != nil {
				logger.Println("❌ Something went wrong while Deploying: ", err)
				fmt.Println(string(deployOutput))
				database.UpdateEventToExistingBuild(buildDetails.BuildId, map[string]map[string]interface{}{"DEPLOY_FAILED": {"timestamp": time.Now(), "reason": string(deployOutput)}})
				return err
			}
			deployedUrl := fmt.Sprintf("http://localhost:%d", port)
			logger.Printf("🥳 Deployed Successfully at port %d. Click here to view the deployed version %s", port, deployedUrl)
			database.UpdateEventToExistingBuild(buildDetails.BuildId, map[string]map[string]interface{}{"DEPLOY_PASSED": {"timestamp": time.Now(), "branded_access_url": deployedUrl, "url": deployedUrl}})
			return nil
		} else {
			logger.Println("🔌 ", port, "is already in use. Trying different port")
		}
	}
	database.UpdateEventToExistingBuild(buildDetails.BuildId, map[string]map[string]interface{}{"DEPLOY_FAILED": {"timestamp": time.Now(), "reason": "All the available 🔌 Ports are used"}})

	return nil
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
