### WebFront
Go Lang project for Automating the build and deployment process of web apps(React).
![Webfront-Miro](https://github.com/Sai7xp/golang-projects/assets/39739036/95af91ca-1490-41d1-b3e4-b472f05e3e7e)

1. `wf-collect-client` Microservice will collect the required details like repo GitHub Url, Build Command to initiate the build process by producing event to kafka
2. `wf-code-builder` this is for processing build events received via Kafka (consumer)
   - Clones the repository from github url received via kafka from wf-collect-client service
   - Generates the Build using the build command
   - Deploys the generated build files

For more details check out the project pdf file in cwd

### Setup
1. Run `docker-compose up -d` to start Apache Kafka & Zookeeper
