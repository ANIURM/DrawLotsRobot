package controller

import (
	"xlab-feishu-robot/app/chat"
	"xlab-feishu-robot/global"
	"xlab-feishu-robot/model"

	"github.com/YasyaKarasu/feishuapi"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

var oldRecordInfo []feishuapi.RecordInfo

func ProjectScheduleReminder(messageevent *chat.MessageEvent) {
	groupID := messageevent.Message.Chat_id
	checkScheduleUpdated(groupID)
}

func checkScheduleUpdated(groupID string) {
	space_id, err := model.QueryKnowledgeSpaceByChat(groupID)
	if err != nil {
		logrus.Warn("[schedule] ", groupID, " get space id fail")
	}
	_, fileToken := getNodeFileToken(space_id, "排期甘特图", "任务进度管理")
	allBitables := global.Feishu.GetAllBitables(fileToken)

	tableInfoList := GetAllTableInfo(allBitables)
	logrus.Debug("[schedule] tableInfoList: ", tableInfoList)

	recordInfoList := GetAllRecordInfo(tableInfoList)
	logrus.Debug("[shchedule] recordInfoList: ", recordInfoList)

	//播报进度数据获取
	var notStarted []string
	var inProgress []string
	var completed []string
	for _, bittable := range allBitables {
		tables := global.Feishu.GetAllTablesInBitable(bittable.AppToken)
		for _, table := range tables {
			records := global.Feishu.GetAllRecordsInTable(bittable.AppToken, table.TableId)
			for _, record := range records {
				if record.Fields["任务状态"] == "未开始" {
					notStarted = append(notStarted, record.Fields["任务名"].(string))
				} else if record.Fields["任务状态"] == "进行中" {
					inProgress = append(inProgress, record.Fields["任务名"].(string))
				} else if record.Fields["任务状态"] == "已完成" {
					completed = append(completed, record.Fields["任务名"].(string))
				}

			}

		}
	}

	user_id, err := model.QueryProjectLeaderByChat(groupID)
	if err != nil {
		logrus.Warn("[schedule] ", groupID, " get project leader fail")
	}
	modified := CheckRecordInfoModified(recordInfoList, oldRecordInfo)

	groupName, err := model.QueryProjectNameByChat(groupID)
	if err != nil {
		logrus.Warn("[schedule] ", groupID, " get project name fail")
	}
	//自动播报
	var msg string
	msg = groupName + ": 当前【任务进度管理】看板任务总数" + string(rune(len(notStarted)+len(inProgress)+len(completed))) + "个，未开始任务" + string(rune(len(notStarted))) + "个，进行中任务" + string(rune(len(inProgress))) + "个，已完成任务" + string(rune(len(completed))) + "个\n"

	if modified {
		//TODO: 获取更新的内容

		updates := ""

		msg = msg + "更新内容：\n" + updates

	} else {
		msg = "请及时更新排期甘特图和任务进度管理，跟进项目进程！ \n" + msg
	}
	link := getlink(space_id)
	msg = msg + "欲了解详细内容，请点击: " + link

	global.Feishu.Send(feishuapi.UserOpenId, user_id, feishuapi.Text, msg)

	oldRecordInfo = recordInfoList

}

func getNodeFileToken(space_id string, topFile string, secFile string) (string, string) {
	var nodeToken, fileToken string
	allNode := global.Feishu.GetAllNodes(space_id)
	for _, node := range allNode {
		if node.Title == topFile {
			allSubNode := global.Feishu.GetAllNodes(space_id, node.NodeToken)
			for _, subNode := range allSubNode {
				if subNode.Title == secFile {
					nodeToken = subNode.NodeToken
					fileToken = subNode.ObjToken
					break
				}
			}
			break
		}
	}
	return nodeToken, fileToken
}

func CheckRecordInfoModified(newRecordInfoList []feishuapi.RecordInfo, oldRecordInfoList []feishuapi.RecordInfo) bool {
	if len(newRecordInfoList) != len(oldRecordInfoList) {
		return true
	}
	for i := 0; i < len(newRecordInfoList); i++ {
		if newRecordInfoList[i].LastModifiedTime != oldRecordInfoList[i].LastModifiedTime {
			return true
		}
	}
	return false
}

func StartProjectScheduleTimer(groupID string, c *cron.Cron) bool {

	// every two days at 9:00
	_, err := c.AddFunc("* * 9 1/2 * *", func() {
		checkScheduleUpdated(groupID)
	})

	if err != nil {
		logrus.Error("[timer] ", groupID, " add project schedule timer fail")
		logrus.Error(err)
		return true
	}

	return false
}

func getlink(spaceId string) string {
	var msg string
	var titles []string
	titles = append(titles, "排期甘特图", "任务进度管理")

	nodes := global.Feishu.GetAllNodes(spaceId)
	for _, value := range nodes {
		if in(value.Title, titles) {
			msg = msg + Url.UrlHead + value.NodeToken + " \n"
		}
		if value.HasChild {
			n := global.Feishu.GetAllNodes(spaceId, value.NodeToken)
			for _, v := range n {
				if in(v.Title, titles) {
					msg = msg + Url.UrlHead + v.NodeToken + "\n"
				}
			}
		}
	}
	return msg
}
