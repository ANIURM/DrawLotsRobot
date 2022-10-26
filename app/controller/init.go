package controller

import (
	"xlab-feishu-robot/app/chat"
	"xlab-feishu-robot/app/dispatcher"
	"xlab-feishu-robot/global"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
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

func InitDebugSpace() {
	global.Rob.GroupSpace["oc_01b58f911445bb053d2d34f2a5546243"] = "7145117180906979330"
	global.Rob.GroupOwner["oc_01b58f911445bb053d2d34f2a5546243"] = "65631d22"
}

func StartGroupTimer(chatID string){
	logrus.Info("[timer] start timer")

	c := cron.New(cron.WithSeconds())
	global.Timer.GTimers[chatID] = c

	StartReviewMeetingTimer(chatID, c)
	StartProjectScheduleTimer(chatID, c)

	c.Start()
	logrus.Info("[timer] start group timer success")
}

func EndGroupTimer(chatID string){
	global.Timer.GTimers[chatID].Stop()
	delete(global.Timer.GTimers, chatID)
}

func startTestTimer(chatID string, c *cron.Cron){
	logrus.Info("[timer] add TestTimer")

	c.AddFunc("* * * * * *", func(){
		logrus.Info("[timer] TestTimer")
		global.Cli.Send("chat_id", chatID, "text", "test")
	})
}