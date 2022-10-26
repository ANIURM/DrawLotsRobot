package controller

import (
	"xlab-feishu-robot/app/chat"
	"xlab-feishu-robot/app/dispatcher"
	"xlab-feishu-robot/global"
)

func InitEvent() {
	dispatcher.RegisterListener(chat.Receive, "im.message.receive_v1")
	InitMessageBind()
	InitDebugSpace()
}

func InitMessageBind() {
	chat.GroupMessageRegister(ReviewMeetingMessage, "复盘")
	chat.GroupMessageRegister(ProjectOver, "结项")
	chat.GroupMessageRegister(ProjectScheduleReminder, "进度更新状态")
	chat.GroupMessageRegister(GetProjectSchedule, "进度获取")
	chat.GroupMessageRegister(GroupHelpMenu, "help")
	chat.P2pMessageRegister(p2pHelpMenu, "help")
}

func StartGroupTimer(chatID string){
	StartReviewMeetingTimer(chatID)
	StartProjectScheduleTimer(chatID)
}

func InitDebugSpace() {
	global.Rob.GroupSpace["oc_01b58f911445bb053d2d34f2a5546243"] = "7145117180906979330"
	global.Rob.GroupOwner["oc_01b58f911445bb053d2d34f2a5546243"] = "65631d22"
}

func EndGroupTimer(chatID string){
}

