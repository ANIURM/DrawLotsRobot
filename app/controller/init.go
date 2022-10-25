package controller

import (
	"xlab-feishu-robot/app/chat"
	"xlab-feishu-robot/app/dispatcher"
)

func InitEvent() {
	// register your handlers here
	dispatcher.RegisterListener(chat.Receive, "im.message.receive_v1")
	InitMessageBind()
}

func InitMessageBind() {
	// register your handlers here
	chat.GroupMessageRegister(ReviewMeetingMessage, "复盘")
	chat.GroupMessageRegister(ProjectOver, "结项")
}

func StartGroupTimer(chatID string){
	StartReviewMeetingTimer(chatID)
}