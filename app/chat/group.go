package chat

import (
	"strings"
	"xlab-feishu-robot/app/controller"
	"xlab-feishu-robot/model"

	"github.com/sirupsen/logrus"
)

var groupMessageMap = make(map[string]messageHandler)

type AuthHandler func(messageevent *model.MessageEvent) model.Privileges

var AuthMap = make(map[string]AuthHandler)

func group(messageevent *model.MessageEvent) {
	switch strings.ToUpper(messageevent.Message.Message_type) {
	case "TEXT":
		groupTextMessage(messageevent)
	default:
		logrus.WithFields(logrus.Fields{"message type": messageevent.Message.Message_type}).Warn("Receive group message, but this type is not supported")
	}
}

func groupTextMessage(messageevent *model.MessageEvent) {
	logrus.WithFields(logrus.Fields{"message content": messageevent.Message.Content}).Info("Receive group TEXT message")

	controller.DrawLotsRobot(messageevent)
}
