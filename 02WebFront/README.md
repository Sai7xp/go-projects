Go Lang project for Automating the build and deployment process of web apps(React).

1. `wf-collect-client` Microservice will collect the required details like repo GitHub Url, Build Command to initiate the build process by producing event to kafka
2. `wf-code-builder` this is for processing build events received via Kafka (consumer)
   - Clone the repository from github url received via kafka from wf-collect-client service
   - Generate the Build using the provided build command
   - Deploy the generated build files

For more details check out the project pdf file in cwd

### Setup
1. Run `docker-compose up -d` to start Apache Kafka & Zookeeper
