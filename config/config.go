package config

import (
	"xlab-feishu-robot/app/controller"
	"xlab-feishu-robot/db"

	"xlab-feishu-robot/app/chat"

	"github.com/YasyaKarasu/feishuapi"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Feishu feishuapi.Config
	Server struct {
		Port int

		// add your configuration fields here
		ExampleField1 string
	}

	// add your configuration fields here
	ExampleField2 struct {
		ExampleField3 int
	}
	LeaderGroup struct {
		chat_id string
	}
	DevGroup struct {
		chat_id string
	}
	FeishuProjectFormPath struct {
		AppTokenForProjectCreat string
		TableIdForProjectCreat  string
		AppTokenForMeeting      string
		TableIdForMeeting       string
	}
	TemplateDocs struct {
		SpaceId         string
		ParentNodeToken string
	}
	UrlStrings struct {
		UrlHead                  string
		UrlForProjectCreate      string
		UrlForGetUserAccessToken string
		UrlForMeeting            string
	}
	Database struct {
		Host       string
		Port       string
		User       string
		Password   string
		AuthSource string
	}
}

var C Config

func ReadConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./")
	viper.AddConfigPath("/etc/xlab-project-robot")

	if err := viper.ReadInConfig(); err != nil {
		logrus.Panic(err)
	}

	if err := viper.Unmarshal(&C); err != nil {
		logrus.Error("Failed to unmarshal config")
	}

	logrus.Info("Configuration file loaded")
}

func SetupFeishuApiClient(cli *feishuapi.AppClient) {
	cli.Conf = C.Feishu
	controller.P = C.FeishuProjectFormPath
	controller.T = C.TemplateDocs
	controller.Url = C.UrlStrings
	chat.LeaderGroupID = C.LeaderGroup.chat_id
	chat.DevGroupID = C.DevGroup.chat_id
}

func SetupDatabase() {
	db.Conf = db.Config{Url: "mongodb://" + C.Database.User + ":" + C.Database.Password + "@" + C.Database.Host + ":" +C.Database.Port + "/" + C.Database.AuthSource}
}
