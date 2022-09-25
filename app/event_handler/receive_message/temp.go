package receiveMessage

import (
	"xlab-feishu-robot/pkg/global"
	"github.com/sirupsen/logrus"
	_ "github.com/YasyaKarasu/feishuapi"
	"time"
)

func init(){
	groupMessageRegister(ReviewMeetingTimer, "复盘")
}

type TableInfo struct {
	AppToken string
	TableID string `json:"table_id"`
	Revision int `json:"revision"`
	Name string `json:"name"`
}

func NewTableInfoWithToken(apptoken string, data map[string]interface{}) *TableInfo{
	return &TableInfo{
		AppToken: apptoken,
		TableID: data["table_id"].(string),
		Revision: int(data["revision"].(float64)),
		Name: data["name"].(string),
	}
}

type RecordInfo struct {
	AppToken string
	TableID string `json:"table_id"`
	RecordID string `json:"record_id"`
	Fields map[string]interface{}
}

func NewRecordInfoWithTokenID(apptoken string, table_id string, data map[string]interface{}) *RecordInfo{
	return &RecordInfo{
		AppToken: apptoken,
		TableID: table_id,
		RecordID: data["record_id"].(string),
		Fields: data["fields"].(map[string]interface{}),
	}
}

type FieldInfo struct {
	Status string
	StartTime float64
	EndTime float64
}

func NewFieldInfo(data map[string]interface{}) *FieldInfo{
	return &FieldInfo{
		Status: data["任务状态"].(string),
		StartTime: data["开始日期"].(float64),
		EndTime: data["截止日期"].(float64),
	}
}

func ReviewMeetingTimer(messageevent *MessageEvent) {
	logrus.Info("review meeting timer")

	//todo: space id
	space_id := "7145117180906979330"
	allNode := global.Cli.GetAllNodes(space_id)
	topFile := "排期甘特图"
	secFile := "任务进度管理"
	var fileToken string

	for _, node := range allNode {
		if node.Title == topFile {
			allSubNode := global.Cli.GetAllNodes(space_id, node.NodeToken)
			for _, subNode := range allSubNode {
				if subNode.Title == secFile {
					fileToken = subNode.ObjToken
					break
				}
			}
			break
		}
	}

	allBitables := global.Cli.GetAllBitables(fileToken)

	//get all the tables
	var tableInfoList []TableInfo
	for _,bitable := range allBitables {
		AppToken := bitable.AppToken
		method := "GET"
		path := "open-apis/bitable/v1/apps/" + AppToken +"/tables"
		query := map[string]string{}
		header := map[string]string{}
		body := map[string]string{}
		resp := global.Cli.Request(method, path, query, header, body)
		tablesRaw := resp["items"].([]interface{})
		for _, tableRaw := range tablesRaw {
			tableInfoList = append(tableInfoList, *NewTableInfoWithToken(AppToken, tableRaw.(map[string]interface{})))
		}
	}

	logrus.Debug("[review meeting] table info list: ", tableInfoList)

	//get all the records
	var recordInfoList []RecordInfo
	for _,table := range tableInfoList {
		AppToken := table.AppToken
		method := "GET"
		path := "open-apis/bitable/v1/apps/" + AppToken +"/tables/" + table.TableID + "/records"
		query := map[string]string{}
		header := map[string]string{}
		body := map[string]string{}
		resp := global.Cli.Request(method, path, query, header, body)
		recordsRaw := resp["items"].([]interface{})
		for _, recordRaw := range recordsRaw {
			recordInfoList = append(recordInfoList, *NewRecordInfoWithTokenID(AppToken, table.TableID, recordRaw.(map[string]interface{})))
		}
	}

	logrus.Debug("[review meeting] record info list: ", recordInfoList)

	//get the "复盘会" records
	var reviewMeeting []FieldInfo
	for _, record := range recordInfoList {
		if name, exist := record.Fields["任务名"]; exist {
			if name == "复盘会"{
				reviewMeeting = append(reviewMeeting, *NewFieldInfo(record.Fields))
				break
			}
		}
	}

	logrus.Debug("review meeting: ", reviewMeeting)

	//check if the review meeting is tomorrow
	var nearby bool
	now := float64(time.Now().Unix()*1000)
	for _, meeting := range reviewMeeting{
		diff := meeting.EndTime - now
		if diff < 86400000 && diff > 0 {
			nearby = true
		}
	}
	if nearby {
		global.Cli.Send("chat_id",messageevent.Message.Chat_id,"text","review meeting is nearby")
	}
	
}