package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var CLIENT *mongo.Client
var MONGO_DB *mongo.Database
var COLLECTION *mongo.Collection
var MONGO_LOGIN string
var MONGO_PASS string
var MONGO_URL string
var MONGO_DB_NAME string
var MONGO_TYPE_CONNECT string

func ConnectMongoDB() {

	options := options.Client().ApplyURI(fmt.Sprintf("%s://%s:%s@%s/", MONGO_TYPE_CONNECT, MONGO_LOGIN, MONGO_PASS, MONGO_URL))

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
	}
}

func DisconnectMongoDB() {
	if CLIENT != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := CLIENT.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}

		log.Println("Disconnected from MongoDB")
	} else {
		log.Println("MongoDB client is not initialized")
	}
}

func InsertOneDbCommit(data CommitData, collectionName string) bool {
	var result bson.M
	getError := getCollection(collectionName).FindOne(context.TODO(), bson.M{"id": data.ID}).Decode(&result)
	if getError != nil {
		res, err := getCollection(collectionName).InsertOne(context.TODO(), data)
		if err != nil {
			log.Fatal(err)
			return false
		}
		log.Println("Inserted ID:", res)
		return true
	}
	return false
}

func getCollection(collectionName string) *mongo.Collection {
	return MONGO_DB.Collection(collectionName)
}

/*
В целом проверка не нужно и этот функционал надо вывести, пока оставлю. Если дого не возвращусь к нему, значит надо будет удалить.
*/
// func checkCollection(collectionName string) bool {
// 	collectionNames, err := MONGO_DB.ListCollectionNames(context.Background(), bson.D{})
// 	if err != nil {
// 		log.Printf("Error checking collection existence: %v", err)
// 		return false
// 	}

// 	for _, name := range collectionNames {
// 		if name == collectionName {
// 			return true
// 		}
// 	}

// 	createCollection(collectionName)
// 	return false
// }
/*
Использовался совместо с функцией checkCollection, то же самое, если долго к ней не вернусь то удалить
*/
// func createCollection(collectionName string) {
// 	if err := MONGO_DB.CreateCollection(context.Background(), collectionName); err != nil {
// 		log.Printf("Error creating collection: %v", err)
// 	} else {
// 		log.Printf("Collection '%s' created successfully", collectionName)
// 	}
// }
