package chat

import (
	"strings"
	"xlab-feishu-robot/app/controller"
	"xlab-feishu-robot/global"
	"xlab-feishu-robot/model"

	"github.com/YasyaKarasu/feishuapi"
	"github.com/sirupsen/logrus"
)

var groupMessageMap = make(map[string]messageHandler)

func group(messageevent *model.MessageEvent) {
	switch strings.ToUpper(messageevent.Message.Message_type) {
	case "TEXT":
		groupTextMessage(messageevent)
	default:
		logrus.WithFields(logrus.Fields{"message type": messageevent.Message.Message_type}).Warn("Receive group message, but this type is not supported")
	}
}

func groupTextMessage(messageevent *model.MessageEvent) {
	// get the pure text message, without @xxx
	messageevent.Message.Content = strings.TrimSuffix(strings.TrimPrefix(messageevent.Message.Content, "{\"text\":\""), "\"}")
	messageevent.Message.Content = messageevent.Message.Content[strings.Index(messageevent.Message.Content, " ")+1:]
	logrus.WithFields(logrus.Fields{"message content": messageevent.Message.Content}).Info("Receive group TEXT message")

	if handler, exists := groupMessageMap[messageevent.Message.Content]; exists {
		privileges := controller.Authenticate(messageevent)
		if privileges == model.Other {
			return
		} else if privileges == model.ProductManagerGroupMembers {
			handler(messageevent)
			return
		} else if privileges == model.ProjectGroupLeader {
			handler(messageevent)
			return
		}
	} else {
		logrus.Error("Group message failed to find event handler: ", messageevent.Message.Content)
		global.Feishu.Send(feishuapi.GroupChatId, messageevent.Message.Chat_id, "text", "关键词"+" ["+messageevent.Message.Content+"] "+"未定义！")
		return
	}
}

func GroupMessageRegister(f messageHandler, s string) {

	if _, isEventExist := groupMessageMap[s]; isEventExist {
		logrus.Warning("Double declaration of group message handler: ", s)
	}
	groupMessageMap[s] = f
}
