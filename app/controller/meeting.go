package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"xlab-feishu-robot/app/chat"
	"xlab-feishu-robot/global"
)

var (
	MyMeeting NewMeeting
)

type NewMeeting struct {
	Code int64 `json:"code"`
	Data struct {
		Record struct {
			Fields struct {
				MeetingName     string                `json:"会议名称"`
				MeetingCategory string                `json:"类别"`
				Moderators      []ParticipatingMember `json:"主持人"`
				Data            int64                 `json:"日期"`
				AttachmentImg   struct {
					Url string `json:"url"`
				} `json:"附件（图片版）"`
				AttachmentText string `json:"附件（网址、文字版）"`
			} `json:"fields"`
			ID       string `json:"id"`
			RecordID string `json:"record_id"`
		} `json:"record"`
	} `json:"data"`
	Msg string `json:"msg"`
}

func ReadMeetingForm(c *gin.Context) {
	resp, _ := c.GetRawData()
	temp := make(map[string]string)
	json.Unmarshal(resp, &temp)
	recordId := temp["record_id"]
	data := global.Cli.GetRecordInByte(P.AppTokenForMeeting, P.TableIdForMeeting, recordId)
	err := json.Unmarshal(data, &MyMeeting)
	if err != nil {
		panic(err)
	}

}

func MeetingForm(event *chat.MessageEvent) {
	msg := "请填写下方的会议问卷：\n" + Url.UrlForMeeting
	global.Cli.Send("chat_id", event.Message.Chat_id, "text", msg)
}
