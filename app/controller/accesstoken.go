package controller

import (
	"github.com/gin-gonic/gin"
	"xlab-feishu-robot/pkg/global"
)

var (
	UserAccessToken string
)

func GetUserAccessToken(c *gin.Context) {
	code, status := c.GetQuery("code")
	if !status {
		panic("The param 'code' was not obtained")
	}
	UserAccessToken = global.Cli.GetUserAccessToken(code).Access_token
	//c.String(200, UserAccessToken)
}
