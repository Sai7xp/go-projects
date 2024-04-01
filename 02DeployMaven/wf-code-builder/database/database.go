/*
* Created on 09 March 2024
* @author Sai Sumanth
 */

package database

import (
	"code-builder-service/models"
	"context"
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

func Init() {
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

func AddNewBuild(newBuildDetails models.BuildRequestDetails) error {
	_, err := buildsCollection.InsertOne(context.TODO(), newBuildDetails)
	if err != nil {
		return err
	}
	fmt.Println("Add New Build details to database with build Id - ", newBuildDetails.BuildId)
	return nil
}

func UpdateEventToExistingBuild(buildId string, newEventData map[string]map[string]interface{}) error {
	// bsonData, err := bson.Marshal(newEventData)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	var firstKey string
	for key := range newEventData {
		firstKey = key
		break // Exit loop after first iteration
	}
	nestedFilePath := fmt.Sprintf("events.%s", firstKey)

	filter := bson.M{"build_id": buildId}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: nestedFilePath, Value: newEventData[firstKey]}}}}

	res := buildsCollection.FindOneAndUpdate(context.TODO(), filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After))
	if res.Err() != nil {
		fmt.Println(res.Err())
		return res.Err()
	}
	fmt.Println("Added new event", firstKey, "to build:", buildId, "successfully!")
	return nil
}
