package controller

import (
	"xlab-feishu-robot/global"

	"github.com/gin-gonic/gin"
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
	c.String(200, "鉴权成功，您现在可以返回，继续您的操作")
}