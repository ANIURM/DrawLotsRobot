package controller

import (
	"time"
	"xlab-feishu-robot/global"
	"xlab-feishu-robot/app/chat"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"github.com/YasyaKarasu/feishuapi"
	"xlab-feishu-robot/global/rob"
)

func ReviewMeetingMessage(messageevent *chat.MessageEvent){
	chatID := messageevent.Message.Chat_id
	haveReviewMeeting := CheckReviewMeeting(chatID)
	if(haveReviewMeeting == 0){
		global.Cli.Send("chat_id", messageevent.Message.Chat_id, "text", "近期没有复盘会，安心啦")
	}
}

type TableInfo struct {
	AppToken string
	TableID  string `json:"table_id"`
	Revision int    `json:"revision"`
	Name     string `json:"name"`
}

func NewTableInfoWithToken(apptoken string, data map[string]interface{}) *TableInfo {
	return &TableInfo{
		AppToken: apptoken,
		TableID:  data["table_id"].(string),
		Revision: int(data["revision"].(float64)),
		Name:     data["name"].(string),
	}
}

type RecordInfo struct {
	AppToken string
	TableID  string `json:"table_id"`
	RecordID string `json:"record_id"`
	Last_modified_time float64 `json:"last_modified_time"`
	Fields   map[string]interface{}
}

func NewRecordInfoWithTokenID(apptoken string, table_id string, data map[string]interface{}) *RecordInfo {
	return &RecordInfo{
		AppToken: apptoken,
		TableID:  table_id,
		RecordID: data["record_id"].(string),
		Last_modified_time: data["last_modified_time"].(float64),
		Fields:   data["fields"].(map[string]interface{}),
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

func CheckReviewMeeting(chatID string) int{

	space_id := rob.Rob.GetGroupSpace(chatID)
	_,fileToken := getNodeFileToken(space_id, "排期甘特图", "任务进度管理")
	allBitables := global.Cli.GetAllBitables(fileToken)

	var tableInfoList []TableInfo
	tableInfoList = GetAllTableInfo(allBitables)
	logrus.Debug("[review meeting] table info list: ", tableInfoList)

	var recordInfoList []RecordInfo
	recordInfoList = GetAllRecordInfo(tableInfoList)
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
		global.Cli.Send("chat_id", chatID, "text", "review meeting is nearby")
		return 1
	}else{
		return 0
	}
}

func GetAllTableInfo(bitableInfoList []feishuapi.BitableInfo) []TableInfo{
	var tableInfoList []TableInfo
	for _, bitable := range bitableInfoList {
		AppToken := bitable.AppToken
		method := "GET"
		path := "open-apis/bitable/v1/apps/" + AppToken + "/tables"
		query := map[string]string{}
		header := map[string]string{}
		body := map[string]string{}
		resp := global.Cli.Request(method, path, query, header, body)
		tablesRaw := resp["items"].([]interface{})
		for _, tableRaw := range tablesRaw {
			tableInfoList = append(tableInfoList, *NewTableInfoWithToken(AppToken, tableRaw.(map[string]interface{})))
		}
	}
	return tableInfoList
}

func GetAllRecordInfo(tableInfoList []TableInfo) []RecordInfo{
	var recordInfoList []RecordInfo
	for _, table := range tableInfoList {
		AppToken := table.AppToken
		method := "GET"
		path := "open-apis/bitable/v1/apps/" + AppToken + "/tables/" + table.TableID + "/records"
		query := map[string]string{} 
		query["automatic_fields"] = "true"
		header := map[string]string{}
		body := map[string]string{}
		resp := global.Cli.Request(method, path, query, header, body)
		recordsRaw := resp["items"].([]interface{})
		for _, recordRaw := range recordsRaw {
			recordInfoList = append(recordInfoList, *NewRecordInfoWithTokenID(AppToken, table.TableID, recordRaw.(map[string]interface{})))
		}
	}
	return recordInfoList
}

// chatID is the groupID
func StartReviewMeetingTimer(chatID string, c *cron.Cron){
	logrus.Info("[timer] ", chatID," add review meeting timer")

	c.AddFunc("0 0 18 * * *", func() {
		CheckReviewMeeting(chatID)
	})
}
