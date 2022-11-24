package controller

import (
	"xlab-feishu-robot/app/chat"
	"xlab-feishu-robot/global"

	"github.com/YasyaKarasu/feishuapi"
)

func p2pHelpMenu(messageevent *chat.MessageEvent) {
	open_id := messageevent.Sender.Sender_id.Open_id
	user_id := messageevent.Sender.Sender_id.User_id
	global.Cli.Send(feishuapi.UserOpenId, open_id, feishuapi.Text, "your user_id is: "+user_id+" , and your open_id is: "+open_id)
}

func GroupHelpMenu(messageevent *chat.MessageEvent) {
	groupID := messageevent.Message.Chat_id
	global.Cli.Send(feishuapi.GroupChatId, groupID, feishuapi.Text, "your groupID is: "+groupID)
}
