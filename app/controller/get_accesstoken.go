package controller

import (
	"xlab-feishu-robot/global"

	"github.com/YasyaKarasu/feishuapi"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	//UserAccessToken string
	UserAccessToken map[string]string
)

func GetUserAccessToken(c *gin.Context) {
	// FIXME
	UserAccessToken = make(map[string]string)
	code, status := c.GetQuery("code")
	if !status {
		panic("The param 'code' was not obtained")
		// how dare you PANIC?
	}
	returnValue := global.Feishu.GetUserAccessToken(code)
	// if == nil ?
	//UserAccessToken = returnValue.Access_token
	UserAccessToken[returnValue.Open_id] = returnValue.Access_token
	TokenUserID = returnValue.User_id
	logrus.Info("UserAccessToken: ", UserAccessToken[returnValue.Open_id])
	logrus.Info("TokenUserID: ", TokenUserID)
	// c.String(200, "鉴权成功，您现在可以返回，继续您的操作")

	//TODO: check if the following code is right
	c.String(200, ` <!DOCTYPE html>
	<html>
	<head>
	</head>
	<body>
		<script>
			window.location.href="about:blank";
			window.close();
		</script>
	</body>
	</html>`)

	global.Feishu.MessageSend(feishuapi.UserUserId, TokenUserID, feishuapi.Text, "鉴权成功，您现在可以返回，继续您的操作")
}
