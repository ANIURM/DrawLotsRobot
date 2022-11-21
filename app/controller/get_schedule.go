package controller

import (
	"xlab-feishu-robot/app/chat"
	"xlab-feishu-robot/global"
	"xlab-feishu-robot/global/robot"

	"github.com/sirupsen/logrus"
)

func GetProjectSchedule(messageevent *chat.MessageEvent) {
	groupID := messageevent.Message.Chat_id
	space_id, ok := robot.Robot.GetGroupSpace(groupID)
	if !ok {
		logrus.WithField("Group ID", groupID).Error("Group space not found")
		return
	}
	nodeToken, fileToken := getNodeFileToken(space_id, "排期甘特图", "任务进度管理")

	logrus.Debug("nodeToken: ", nodeToken, " fileToken: ", fileToken)

	user_id, ok := robot.Robot.GetGroupOwner(groupID)
	if !ok {
		logrus.WithField("Group ID", groupID).Error("Group owner not found")
		return
	}
	fileLink := Url.UrlHead + nodeToken

	groupInfo := global.Cli.GetGroupInfo(groupID)
	groupName := groupInfo.Name
	global.Cli.Send("open_id", user_id, "text", groupName+" 任务进度： "+fileLink)
}
