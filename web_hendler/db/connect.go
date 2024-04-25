package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var CLIENT *mongo.Client
var MONGO_DB *mongo.Database
var COLLECTION *mongo.Collection

func ConnectMongoDB() {
	options := options.Client().ApplyURI("mongodb://admin:my_password@localhost:27017/")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Connected to MongoDB!")
		CLIENT = client
		MONGO_DB = client.Database("local")
		COLLECTION = MONGO_DB.Collection("commits")
	}

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		if err := client.Disconnect(ctx); err != nil {
			log.Fatalf("Error disconnecting from MongoDB: %v", err)
		}
		log.Println("Disconnected from MongoDB")
	}()
}

func InsertOneDbCommit(data CommitData, collectionName string) bool {

	checkCollection(collectionName)

	res, err := COLLECTION.InsertOne(context.TODO(), data)
	if err != nil {
		log.Fatal(err)
		return false
	} else {
		log.Println("Inserted ID:", res)
		return true
	}

}

func checkCollection(collectionName string) bool {
	collectionNames, err := MONGO_DB.ListCollectionNames(context.Background(), nil)
	if err != nil {
		log.Printf("Error checking collection existence: %v", err)
		return false
	}

	for _, name := range collectionNames {
		if name == collectionName {
			return true
		}
	}

	createCollection(collectionName)
	return false
}

func createCollection(collectionName string) {
	if err := MONGO_DB.CreateCollection(context.Background(), collectionName); err != nil {
		log.Printf("Error creating collection: %v", err)
	} else {
		log.Printf("Collection '%s' created successfully", collectionName)
	}
}
