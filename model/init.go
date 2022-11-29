package model

import (
	"context"
	"xlab-feishu-robot/db"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Config struct {
	Host       string
	Port       string
	User       string
	Password   string
	AuthSource string
}

var (
	client  *mongo.Client
	project *mongo.Collection
)

func InitDatabase() {
	client = db.Connect()
	project = client.Database("xlabFeishuRobot").Collection("project")
	task = client.Database("xlabFeishuRobot").Collection("task")

	// create index on ChatID
	name, err := project.Indexes().CreateOne(context.TODO(), mongo.IndexModel{Keys: bson.D{{Key: "ProjectChat", Value: 1}}})
	if err != nil {
		logrus.Error(err)
	}
	logrus.Info("Created index on ProjectChat: ", name)
}
