package model

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var employee, group, privilege, project *mongo.Collection

func Init() {
	var err error

	clientOptions := options.Client().ApplyURI("mongodb://test:123456@101.34.188.18:27017/test")

	client, err = mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	employee = client.Database("test").Collection("employee")
	group = client.Database("test").Collection("group")
	privilege = client.Database("test").Collection("privilege")
	project = client.Database("test").Collection("project")
}
