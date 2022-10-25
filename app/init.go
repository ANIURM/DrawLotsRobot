package app

import (
	"xlab-feishu-robot/app/controller"
	"xlab-feishu-robot/app/dispatcher"

	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine) {
	controller.InitEvent()
}

func Register(r *gin.Engine) {
	// register your controllers here
	// example
	//r.POST("/api/example", controller.Example)
	r.GET("/api/getUserAccessToken", controller.GetUserAccessToken)
	r.POST("/api/project", controller.InitProject)

	// DO NOT CHANGE LINES BELOW
	// register dispatcher
	r.POST("/feiShu/Event", dispatcher.Dispatcher)
}
