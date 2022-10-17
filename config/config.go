package config

import (
	"github.com/YasyaKarasu/feishuapi"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"xlab-feishu-robot/app/controller"
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
		AppToken string
		TableId  string
	}
	TemplateDocs struct {
		SpaceId         string
		ParentNodeToken string
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
}
