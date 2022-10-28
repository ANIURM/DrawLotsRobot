package config

import (
	"xlab-feishu-robot/app/controller"

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
}

var C Config

func ReadConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./")

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
}
