package controller

import (
	"xlab-feishu-robot/global"
	"xlab-feishu-robot/model"

	"github.com/YasyaKarasu/feishuapi"
	"github.com/sirupsen/logrus"
)

func GetProjectSchedule(messageevent *model.MessageEvent) {
	groupID := messageevent.Message.Chat_id
	space_id, err := model.QueryKnowledgeSpaceByChat(groupID)
	if err != nil {
		logrus.Warn("[schedule] ", groupID, " get space id fail")
	}
	nodeToken, fileToken := getNodeFileToken(space_id, "排期甘特图", "任务进度管理")

	logrus.Debug("nodeToken: ", nodeToken, " fileToken: ", fileToken)

	user_id, err := model.QueryProjectLeaderByChat(groupID)
	if err != nil {
		logrus.Warn("[schedule] ", groupID, " get project leader fail")
		return
	}
	fileLink := Url.UrlHead + nodeToken

	groupName, err := model.QueryProjectNameByChat(groupID)
	if err != nil {
		logrus.Warn("[schedule] ", groupID, " get project name fail")
	}
	global.Feishu.MessageSend(feishuapi.UserOpenId, user_id, feishuapi.Text, groupName+" 任务进度： "+fileLink)
}
