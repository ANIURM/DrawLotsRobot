package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"xlab-feishu-robot/app"
	"xlab-feishu-robot/config"
	"xlab-feishu-robot/docs"
	"xlab-feishu-robot/global"
)

func main() {

	config.ReadConfig()

	// log
	config.SetupLogrus()
	logrus.Info("Robot starts up")

	// debug 模式
	// logrus.SetLevel(logrus.DebugLevel)

	// feishu api client
	config.SetupFeishuApiClient(&global.Feishu)
	global.Feishu.StartTokenTimer()

	// robot server
	r := gin.Default()
	app.Init(r)

	// api docs by swagger
	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.Run(":" + fmt.Sprint(config.C.Server.Port))
}
