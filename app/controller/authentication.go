package controller

import (
	"github.com/YasyaKarasu/feishuapi"
	"github.com/sirupsen/logrus"
	"xlab-feishu-robot/global"
	"xlab-feishu-robot/model"
)

var (
	LeaderGroupID string
	DevGroupID    string
)

type messageHandler func(event *model.MessageEvent)

var groupMessageMap = make(map[string]messageHandler)

func Authenticate(messageevent *model.MessageEvent) model.Privileges {
	// "立项"是产品经理群成员权限
	if messageevent.Message.Content == "立项" {
		chat_id := messageevent.Message.Chat_id
		if chat_id != LeaderGroupID && chat_id != DevGroupID {
			global.Feishu.Send(feishuapi.GroupChatId, messageevent.Message.Chat_id, feishuapi.Text, "抱歉，您暂时没有权限进行 ["+messageevent.Message.Content+"] 操作！")
			return model.Other
		}
		return model.ProductManagerGroupMembers
	} else {
		// 项目群群主权限
		chat_id := messageevent.Message.Chat_id
		owner_id, err := model.QueryProjectLeaderByChat(chat_id)
		if err != nil {
			logrus.Error("[authenticate] ", chat_id, " get project leader id fail")
		}
		if chat_id != DevGroupID && messageevent.Sender.Sender_id.Open_id != owner_id {
			global.Feishu.Send(feishuapi.GroupChatId, messageevent.Message.Chat_id, feishuapi.Text, "抱歉，您暂时没有权限进行 ["+messageevent.Message.Content+"] 操作！")
			return model.Other
		} else {
			return model.ProjectGroupLeader
		}
	}

}
