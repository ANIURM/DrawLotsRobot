package db

import (
	"context"
	"log"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Url string
}

var Conf Config

func Connect() *mongo.Client{
	var err error
	clientOpts := options.Client().ApplyURI(Conf.Url)
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	logrus.Info("Connected to MongoDB!")
	return client
}