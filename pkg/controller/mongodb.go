package controller

import (
	"context"
	"os"
	"time"

	"vngitSub/pkg/utils"
	"vngitSub/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connectMongo() (*mongo.Client, string, ) {
	logger := utils.ConfigZap()

	mongoAddress := os.Getenv("MONGO_ADDRESS")
	mongoDBName := os.Getenv("MONGO_DBNAME")
	mongoUsername := os.Getenv("MONGO_USER")
	mongoPassword := os.Getenv("MONGO_PASS")

	uri := "mongodb://" + mongoUsername + ":" + mongoPassword + "@" + mongoAddress + "/?retryWrites=true&w=majority"

	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI).SetConnectTimeout(3*time.Second)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)

	if err != nil {
		logger.Errorf("Connecting to MongoDB...FAILED: %s", err)
	}
	// defer func() {
	// 	if err = client.Disconnect(context.TODO()); err != nil {
	// 		log.Panic(err)
	// 	}
	// }()

	return client, mongoDBName
}

func ValidateMongoConnection() {
	logger := utils.ConfigZap()
	client, mongoDBName := connectMongo()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Check connection to MongoDB
	var result bson.M
	if err := client.Database(mongoDBName).RunCommand(ctx, bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		logger.Errorf("Connecting to MongoDB...FAILED: %s", err)
	} else {
		logger.Info("Connecting to MongoDB...OK")
	}
}

func saveState(transID string, image string, oldTag string, newTag string, cluster string, blobName string, time string, status string) {
	logger := utils.ConfigZap()
	client, mongoDBName := connectMongo()

	coll := client.Database(mongoDBName).Collection("status")
	newStatus := model.MessageStatus{Image: image, OldTag: oldTag, NewTag: newTag, Cluster: cluster, BlobName: blobName, Time: time, Status: status, Metadata: "Sent from PubSub version"}

	result, err := coll.InsertOne(context.TODO(), newStatus)
	if err != nil {
		logger.Errorf("[%s] Saving image state to MongoDB...FAILED: %s", transID, err)
	}
	logger.Infof("[%s] Saving image state to MongoDB with ID [%s]...OK", transID, result)
}