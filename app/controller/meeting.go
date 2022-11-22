package controller

import (
	"xlab-feishu-robot/app/chat"
	"xlab-feishu-robot/global"
	"xlab-feishu-robot/model"
)

// var (
// 	MyMeeting NewMeeting
// )

// type NewMeeting struct {
// 	Code int64 `json:"code"`
// 	Data struct {
// 		Record struct {
// 			Fields struct {
// 				MeetingName     string                `json:"会议名称"`
// 				MeetingCategory string                `json:"类别"`
// 				Moderators      []ParticipatingMember `json:"主持人"`
// 				Data            int64                 `json:"日期"`
// 				AttachmentImg   struct {
// 					Url string `json:"url"`
// 				} `json:"附件（图片版）"`
// 				AttachmentText string `json:"附件（网址、文字版）"`
// 			} `json:"fields"`
// 			ID       string `json:"id"`
// 			RecordID string `json:"record_id"`
// 		} `json:"record"`
// 	} `json:"data"`
// 	Msg string `json:"msg"`
// }

// func ReadMeetingForm(c *gin.Context) {
// 	resp, _ := c.GetRawData()
// 	temp := make(map[string]string)
// 	json.Unmarshal(resp, &temp)
// 	recordId := temp["record_id"]
// 	data := global.Cli.GetRecordInByte(P.AppTokenForMeeting, P.TableIdForMeeting, recordId)
// 	err := json.Unmarshal(data, &MyMeeting)
// 	if err != nil {
// 		panic(err)
// 	}
// }

func MeetingForm(messageevent *chat.MessageEvent) {
	groupID := messageevent.Message.Chat_id
	space_id, err := model.GetKnowledgeSpaceByChat(groupID)
	if err != nil {
		return
	}

	nodes := global.Cli.GetAllNodes(space_id)
	var node_token string
	for _, value := range nodes {
		if value.Title == "项目会议" {
			node_token = value.NodeToken
		}
	}
	msg := "请填写下方的会议问卷：\n" + Url.UrlHead + node_token
	global.Cli.Send("chat_id", messageevent.Message.Chat_id, "text", msg)
}
