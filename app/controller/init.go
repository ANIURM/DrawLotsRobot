package controller

import (
	"xlab-feishu-robot/app/chat"
	"xlab-feishu-robot/app/dispatcher"
	"xlab-feishu-robot/global"
	"xlab-feishu-robot/global/robot"
	"xlab-feishu-robot/model"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

func InitEvent() {
	dispatcher.RegisterListener(chat.Receive, "im.message.receive_v1")
	InitMessageBind()
	InitDebugSpace()
	recoverTimer()
	// Debug()
}

func InitMessageBind() {
	// 产品经理群所有人
	chat.GroupMessageRegister(ProjectCreat, "立项")
	// 项目群所有者
	chat.GroupMessageRegister(UpdateMember, "人员变更")
	chat.GroupMessageRegister(ReviewMeetingMessage, "复盘")
	chat.GroupMessageRegister(ProjectOver, "结项")
	chat.GroupMessageRegister(ProjectScheduleReminder, "进度更新状态")
	chat.GroupMessageRegister(GetProjectSchedule, "进度获取")
	chat.GroupMessageRegister(MeetingForm, "会议问卷")
	chat.GroupMessageRegister(GroupHelpMenu, "help")
	// 所有人
	chat.P2pMessageRegister(p2pHelpMenu, "help")
}

func InitDebugSpace() {
	robot.Robot.SetGroupSpace("oc_01b58f911445bb053d2d34f2a5546243", "7145117180906979330")
	robot.Robot.SetGroupOwner("oc_01b58f911445bb053d2d34f2a5546243", "65631d22")
}

func StartGroupTimer(chatID string) {

	c := cron.New(cron.WithSeconds())
	global.Timer.GTimers[chatID] = c

	StartReviewMeetingTimer(chatID, c)
	StartProjectScheduleTimer(chatID, c)

	c.Start()
	logrus.Info("[timer] group [ ", chatID, " ] start group time")
}

func EndGroupTimer(chatID string) {
	global.Timer.GTimers[chatID].Stop()
	delete(global.Timer.GTimers, chatID)
}

func startTestTimer(chatID string, c *cron.Cron) {
	logrus.Info("[timer] add TestTimer")

	c.AddFunc("* * * * * *", func() {
		logrus.Info("[timer] TestTimer")
		global.Cli.Send("chat_id", chatID, "text", "test")
	})
}

func recoverTimer() {
	logrus.Info("[timer] ------------- recovering timer ---------------")
	for k, _ := range robot.Robot.GetGroupSpaceMap() {
		StartGroupTimer(k)
	}
	logrus.Info("[timer] ------------- recover timer finish ---------------")
}

func Debug() {
	model.DeleteRobotStateRecords("testGroup")
	robot.Robot.SetGroupSpace("testGroup", "testSpace")
	robot.Robot.SetGroupOwner("testGroup", "testUser")
	space, ok := robot.Robot.GetGroupSpace("testGroup")
	if !ok {
		logrus.WithField("Group ID", "testGroup").Error("Group space not found")
		return
	}
	user, ok := robot.Robot.GetGroupOwner("testGroup")
	if !ok {
		logrus.WithField("Group ID", "testGroup").Error("Group owner not found")
		return
	}
	if space == "testSpace" && user == "testUser" {
		logrus.Info("[debug] insert test success")
	}
	robot.Robot.SetGroupSpace("testGroup", "testSpace222")
	robot.Robot.SetGroupOwner("testGroup", "testUser222")
	space, ok = robot.Robot.GetGroupSpace("testGroup")
	if !ok {
		logrus.WithField("Group ID", "testGroup").Error("Group space not found")
		return
	}
	user, ok = robot.Robot.GetGroupOwner("testGroup")
	if !ok {
		logrus.WithField("Group ID", "testGroup").Error("Group owner not found")
		return
	}
	if space == "testSpace222" && user == "testUser222" {
		logrus.Info("[debug] update test success")
	}
}
