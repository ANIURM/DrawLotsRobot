package controller

import (
	"time"
	"xlab-feishu-robot/global"
	"xlab-feishu-robot/model"

	"github.com/YasyaKarasu/feishuapi"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

func ReviewMeetingMessage(messageevent *model.MessageEvent) {
	chatID := messageevent.Message.Chat_id
	haveReviewMeeting := CheckReviewMeeting(chatID)
	if haveReviewMeeting == 0 {
		global.Feishu.Send(feishuapi.GroupChatId, messageevent.Message.Chat_id, feishuapi.Text, "近期没有复盘会，安心啦")
	}
}

type FieldInfo struct {
	Status    string
	StartTime float64
	EndTime   float64
}

func NewFieldInfo(data map[string]interface{}) *FieldInfo {
	return &FieldInfo{
		Status:    data["任务状态"].(string),
		StartTime: data["开始日期"].(float64),
		EndTime:   data["截止日期"].(float64),
	}
}

func CheckReviewMeeting(chatID string) int {

	space_id, err := model.QueryKnowledgeSpaceByChat(chatID)
	if err != nil {
		return 0
	}

	_, fileToken := getNodeFileToken(space_id, "排期甘特图", "任务进度管理")
	allBitables := global.Feishu.GetAllBitables(fileToken)

	tableInfoList := GetAllTableInfo(allBitables)
	logrus.Debug("[review meeting] table info list: ", tableInfoList)

	recordInfoList := GetAllRecordInfo(tableInfoList)
	logrus.Debug("[review meeting] record info list: ", recordInfoList)

	//get the "复盘会" records
	var reviewMeeting []FieldInfo
	for _, record := range recordInfoList {
		if name, exist := record.Fields["任务名"]; exist {
			if name == "复盘会" {
				reviewMeeting = append(reviewMeeting, *NewFieldInfo(record.Fields))
				break
			}
		}
	}

	logrus.Debug("review meeting: ", reviewMeeting)

	//check if the review meeting is tomorrow
	var nearby bool
	now := float64(time.Now().Unix() * 1000)
	for _, meeting := range reviewMeeting {
		diff := meeting.EndTime - now
		if diff < 86400000 && diff > 0 {
			nearby = true
		}
	}
	if nearby {
		global.Feishu.Send(feishuapi.GroupChatId, chatID, feishuapi.Text, "项目已到复盘会时间，请查看项目进程，若已完成将按时复盘，若未完成请更改复盘会时间")
		return 1
	} else {
		return 0
	}
}

func GetAllTableInfo(bitableInfoList []feishuapi.BitableInfo) []feishuapi.TableInfo {
	var tableInfoList []feishuapi.TableInfo
	for _, bitable := range bitableInfoList {
		tables := global.Feishu.GetAllTablesInBitable(bitable.AppToken)
		tableInfoList = append(tableInfoList, tables...)
	}
	return tableInfoList
}

func GetAllRecordInfo(tableInfoList []feishuapi.TableInfo) []feishuapi.RecordInfo {
	var recordInfoList []feishuapi.RecordInfo
	for _, table := range tableInfoList {
		records := global.Feishu.GetAllRecordsInTable(table.AppToken, table.TableId)
		recordInfoList = append(recordInfoList, records...)
	}
	return recordInfoList
}

// chatID is the groupID
func StartReviewMeetingTimer(chatID string, c *cron.Cron) bool {

	_, err := c.AddFunc("0 0 18 * * *", func() {
		CheckReviewMeeting(chatID)
	})

	if err != nil {
		logrus.Error("[timer]", chatID, "Review Meeting Timer start error")
		return true
	}

	return false
}
