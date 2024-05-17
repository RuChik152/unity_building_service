package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
	"web_hendler/loger"

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
		log.Fatal("Ошибка подключения к базе данных: ", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Ошибка проверки доступности доступа к БД: ", err)
	} else {
		loger.LogPrint.Package("DB").Log("Connected to MongoDB!")
		CLIENT = client
		MONGO_DB = client.Database(MONGO_DB_NAME)
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
		loger.LogPrint.Package("DB").Log(fmt.Sprint("Inserted ID:", res))
		return true
	}
	return false
}

func GetCommitData(id int, collectionName string) (GetCommit, error) {
	var result bson.M
	err := getCollection(collectionName).FindOne(context.TODO(), bson.M{"id": id}).Decode(&result)
	if err != nil {
		return GetCommit{}, errors.New("ошибка получения данных по комиту из базы данных")
	}

	log.Println("TEST >>>> ", result)

	return GetCommit{
		ID:      int(result["id"].(int32)),
		AUTHOR:  result["author"].(string),
		MESSAGE: result["comment"].(string),
		SHA:     result["sha"].(string),
	}, nil
}

func getCollection(collectionName string) *mongo.Collection {
	return MONGO_DB.Collection(collectionName)
}
