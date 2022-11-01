package model

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Host     string
	Port     string
	User 		string
	Password string
	AuthSource   string
}

var Conf Config
var client *mongo.Client
var employee, group, privilege, project,robot_state *mongo.Collection

func InitDatabase() {
	Connect()

	employee = client.Database("xlabFeishuRobot").Collection("employee")
	group = client.Database("xlabFeishuRobot").Collection("group")
	privilege = client.Database("xlabFeishuRobot").Collection("privilege")
	project = client.Database("xlabFeishuRobot").Collection("project")
	robot_state = client.Database("xlabFeishuRobot").Collection("robot_state")

}

func Connect(){
	var err error
	credential := options.Credential{
		AuthSource: Conf.AuthSource,
		Username:	Conf.User,
		Password:  	Conf.Password,
	}
	clientOpts := options.Client().ApplyURI("mongodb://"+Conf.Host+":"+Conf.Port).SetAuth(credential)
	client, err = mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	logrus.Info("Connected to MongoDB!")
}