package controller

import (
	"xlab-feishu-robot/app/chat"
	"xlab-feishu-robot/global"
)

func p2pHelpMenu(messageevent *chat.MessageEvent) {
	open_id := messageevent.Sender.Sender_id.Open_id
	user_id := messageevent.Sender.Sender_id.User_id
	global.Cli.Send("open_id",open_id, "text", "your user_id is: " + user_id + " , and your open_id is: " + open_id)
}

func GroupHelpMenu(messageevent *chat.MessageEvent) {
	groupID := messageevent.Message.Chat_id
	global.Cli.Send("chat_id", groupID, "text", "your groupID is: " + groupID)
}