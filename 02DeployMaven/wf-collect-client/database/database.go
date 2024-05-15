/*
* Created on 09 March 2024
* @author Sai Sumanth
 */

package database

import (
	"context"
	"errors"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dataBaseName   = "webfront_db"
	collectionName = "builds"
	mongoUri       = "mongodb://localhost:27017"
)

var buildsCollection *mongo.Collection

// Initializes the Mongo Database
func Init() {
	fmt.Println("Connecting to MongoDb...")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUri))
	if err != nil {
		fmt.Println("Error while connecting to Database")
		log.Fatal(err)
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("MongoDB Connection Successfulll!! 🎊")
	buildsCollection = client.Database(dataBaseName).Collection(collectionName)
}

func FetchBuildDetailsFromDb(buildId string) ([]byte, error) {
	if buildsCollection == nil {
		return nil, errors.New("database not donnected")
	}
	filter := bson.D{{Key: "build_id", Value: buildId}}
	result := bson.D{}
	response := buildsCollection.FindOne(context.TODO(), filter)
	if response == nil || response.Err() != nil {
		fmt.Println("Failed to fetch details for given build id - ", buildId)
		return nil, response.Err()
	}
	response.Decode(&result)
	encodedBytes, _ := bson.MarshalExtJSON(result, true, true)
	fmt.Println("Encoded Build Details : ", string(encodedBytes))
	return encodedBytes, nil
}
