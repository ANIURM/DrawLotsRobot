package controller

import (
	"xlab-feishu-robot/app/chat"
	"xlab-feishu-robot/app/dispatcher"
)

func InitEvent() {
	dispatcher.RegisterListener(chat.Receive, "im.message.receive_v1")
}
