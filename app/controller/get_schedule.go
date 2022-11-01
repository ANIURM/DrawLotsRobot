package controller

import (	
	"xlab-feishu-robot/app/chat"
	"xlab-feishu-robot/global"
	"xlab-feishu-robot/global/rob"
	"github.com/sirupsen/logrus"
)

func GetProjectSchedule(messageevent *chat.MessageEvent){
	groupID := messageevent.Message.Chat_id
	space_id := rob.Rob.GetGroupSpace(groupID)
	nodeToken, fileToken := getNodeFileToken(space_id, "排期甘特图", "任务进度管理")
	
	logrus.Debug("nodeToken: ", nodeToken, " fileToken: ", fileToken)

	user_id := rob.Rob.GetGroupOwner(groupID)
	fileLink := "xn4zlkzg4p.feishu.cn/wiki/" + nodeToken
	global.Cli.Send("user_id", user_id, "text", fileLink)
}