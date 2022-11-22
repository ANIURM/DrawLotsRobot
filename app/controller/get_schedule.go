package controller

import (
	"xlab-feishu-robot/app/chat"
	"xlab-feishu-robot/global"
	"xlab-feishu-robot/model"

	"github.com/YasyaKarasu/feishuapi"
	"github.com/sirupsen/logrus"
)

func GetProjectSchedule(messageevent *chat.MessageEvent) {
	groupID := messageevent.Message.Chat_id
	space_id, err := model.GetKnowledgeSpaceByChat(groupID)
	if err != nil {
		logrus.Warn("[schedule] ", groupID, " get space id fail")
	}
	nodeToken, fileToken := getNodeFileToken(space_id, "排期甘特图", "任务进度管理")

	logrus.Debug("nodeToken: ", nodeToken, " fileToken: ", fileToken)

	user_id,err := model.GetProjectLeaderByChat(groupID)
	if err != nil {
		logrus.Warn("[schedule] ", groupID, " get project leader fail")
		return 
	}
	fileLink := Url.UrlHead + nodeToken

	groupInfo := global.Cli.GetGroupInfo(groupID)
	groupName := groupInfo.Name
	global.Cli.Send(feishuapi.UserUserId, user_id, "text", groupName+" 任务进度： "+fileLink)
}
