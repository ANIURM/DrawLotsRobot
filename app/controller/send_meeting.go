package controller

import (
	"xlab-feishu-robot/global"
	"xlab-feishu-robot/model"

	"github.com/YasyaKarasu/feishuapi"
)

func MeetingForm(messageevent *model.MessageEvent) {
	groupID := messageevent.Message.Chat_id
	space_id, err := model.QueryKnowledgeSpaceByChat(groupID)
	if err != nil {
		return
	}

	nodes := global.Feishu.KnowledgeSpaceGetAllNodes(space_id)
	var node_token string
	for _, value := range nodes {
		if value.Title == "项目会议" {
			node_token = value.NodeToken
		}
	}
	msg := "请填写下方的会议问卷：\n" + Url.UrlHead + node_token
	global.Feishu.MessageSend(feishuapi.GroupChatId, messageevent.Message.Chat_id, feishuapi.Text, msg)
}
