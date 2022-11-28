package controller

import (
	"xlab-feishu-robot/global"
	"xlab-feishu-robot/model"

	"github.com/YasyaKarasu/feishuapi"
)

func p2pHelpMenu(messageevent *model.MessageEvent) {
	open_id := messageevent.Sender.Sender_id.Open_id
	user_id := messageevent.Sender.Sender_id.User_id
	global.Feishu.Send(feishuapi.UserOpenId, open_id, feishuapi.Text, "your user_id is: "+user_id+" , and your open_id is: "+open_id)
}

func GroupHelpMenu(messageevent *model.MessageEvent) {
	groupID := messageevent.Message.Chat_id
	global.Feishu.Send(feishuapi.GroupChatId, groupID, feishuapi.Text, "your groupID is: "+groupID)
}
