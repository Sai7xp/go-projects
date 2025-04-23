### Remote Deploy

Automates the build and deployment process of web apps(supports React for now).
Kafka is used for the build queue.
![Webfront-Miro](https://github.com/Sai7xp/golang-projects/assets/39739036/95af91ca-1490-41d1-b3e4-b472f05e3e7e)

1. `wf-collect-client` Microservice will collect the required details like repo public GitHub Url and Build Command to initiate the build process by producing event to kafka
2. `wf-code-builder` this is for processing build events received via Kafka (consumer)
   - Clones the repository from github url received via kafka from wf-collect-client service
   - Generates the Build using the build command
   - Deploys the generated build files
3. MongoDb is being used for logging the events.

For more details about the project check out the pdf file in cwd.

### Tech Stack
- Golang - Two microservices built using Go
- Kafka - Message broker for event-driven communication
- MongoDB - To store the event logs for each build
- Docker Compose - To spin up Kafka and Zookeeper on docker

### Setup

1. Run `docker-compose up -d` to start Apache Kafka & Zookeeper

### Go Lang & Kafka Implementation

- On Kafka, Microservices can publish events
- Microservices can also subscribe to the events(consume)
- Kafka can store events for later retrieval
