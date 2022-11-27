package controller

import (
	"xlab-feishu-robot/app/chat"
	"xlab-feishu-robot/app/dispatcher"
	"xlab-feishu-robot/global"
	"xlab-feishu-robot/model"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

func InitEvent() {
	dispatcher.RegisterListener(chat.Receive, "im.message.receive_v1")
	InitMessageBind()
	// InitDebugSpace()
	// TestDB()
	recoverTimer()
}

func InitMessageBind() {
	//TODO: 鉴权
	chat.GroupMessageRegister(ProjectCreat, "立项")
	chat.GroupMessageRegister(UpdateMember, "人员变更")
	chat.GroupMessageRegister(ReviewMeetingMessage, "复盘")
	chat.GroupMessageRegister(FinishProject, "结项")
	chat.GroupMessageRegister(ProjectScheduleReminder, "进度更新状态")
	chat.GroupMessageRegister(GetProjectSchedule, "进度获取")
	chat.GroupMessageRegister(MeetingForm, "会议问卷")
	chat.GroupMessageRegister(GroupHelpMenu, "help")
	chat.P2pMessageRegister(p2pHelpMenu, "help")
}

// ! 只调用一次，用于初始化 DEBUG 信息
func InitDebugSpace() {
	// lemon test
	var project model.Project
	project.ProjectName = "lemon test"
	project.ProjectType = model.ProjectType("testProjectType")
	project.ProjectLeaderId = "65631d22" // 其实应该用 open_id
	project.ProjectSpace = "7145117180906979330"
	project.ProjectChat = "oc_01b58f911445bb053d2d34f2a5546243"
	project.ProjectStatus = model.ProStatus("testProjectStatus")

	var projectList []model.Project
	projectList = append(projectList, project)
	model.InsertProjectRecords(projectList)

	// timer should be started
	project.ProjectName = "start timer"
	project.ProjectType = model.ProjectType("testProjectType")
	project.ProjectLeaderId = "65631d22" // 其实应该用 open_id
	project.ProjectSpace = "7145117180906979330"
	project.ProjectChat = "start timer"
	project.ProjectStatus = model.Pending

	projectList = nil
	projectList = append(projectList, project)
	model.InsertProjectRecords(projectList)

	// timer should not be started
	project.ProjectName = "not start timer"
	project.ProjectType = model.ProjectType("testProjectType")
	project.ProjectLeaderId = "65631d22" // 其实应该用 open_id
	project.ProjectSpace = "7145117180906979330"
	project.ProjectChat = "not start timer"
	project.ProjectStatus = model.Finished

	projectList = nil
	projectList = append(projectList, project)
	model.InsertProjectRecords(projectList)

	logrus.Info("[debug] init debug space done")
}

func TestDB() {
	chatID := "oc_01b58f911445bb053d2d34f2a5546243"
	leader, err := model.QueryProjectLeaderByChat(chatID)
	if err != nil || leader != "65631d22" {
		logrus.Error("test db failed")
	}
	space, err := model.QueryKnowledgeSpaceByChat(chatID)
	if err != nil || space != "7145117180906979330" {
		logrus.Error("test db failed")
	}
}

func recoverTimer() {
	logrus.Info("[timer] -------------   recovering timer   ---------------")
	ChatStatusMap, err := model.QueryChatStatusMap()
	if err != nil {
		logrus.Error("[timer] get chat status map failed")
		return
	}
	for chatID, status := range ChatStatusMap {
		if status != model.Finished && status != model.Aborted {
			StartGroupTimer(chatID)
		}
	}
}

func StartGroupTimer(chatID string) {
	c := cron.New(cron.WithSeconds())
	global.Timer.GTimers[chatID] = c

	StartReviewMeetingTimer(chatID, c)
	StartProjectScheduleTimer(chatID, c)

	c.Start()

	groupName, err := model.QueryProjectNameByChat(chatID)
	if err != nil {
		logrus.Error("[timer] get project name by chat failed")
		return
	}

	logrus.Info("[timer] group [ ", groupName, " ] start group time finish")
}

func EndGroupTimer(chatID string) {
	global.Timer.GTimers[chatID].Stop()
	delete(global.Timer.GTimers, chatID)
}
