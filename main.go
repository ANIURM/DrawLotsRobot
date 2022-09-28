package main

import (
	"fmt"
	"xlab-feishu-robot/app"
	"xlab-feishu-robot/config"
	"xlab-feishu-robot/docs"
	"xlab-feishu-robot/pkg/global"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	// "xlab-feishu-robot/app/timer"
)

func main() {
	config.ReadConfig()

	// log
	config.SetupLogrus()
	logrus.Info("Robot starts up")

	// feishu api client
	config.SetupFeishuApiClient(&global.Cli)
	global.Cli.StartTokenTimer()

	// robot server
	r := gin.Default()
	app.Init(r)

	// api docs by swagger
	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	//DEBUG:
	//timer.StartReviewMeetingTimer("oc_01b58f911445bb053d2d34f2a5546243")

	r.Run(":" + fmt.Sprint(config.C.Server.Port))
}
