package chat

import (
	"strings"
	"xlab-feishu-robot/global"
	"xlab-feishu-robot/global/robot"

	"github.com/sirupsen/logrus"
)

var (
	LeaderGroupID string
	DevGroupID string
)

var groupMessageMap = make(map[string]messageHandler)

func group(messageevent *MessageEvent) {
	switch strings.ToUpper(messageevent.Message.Message_type) {
	case "TEXT":
		groupTextMessage(messageevent)
	default:
		logrus.WithFields(logrus.Fields{"message type": messageevent.Message.Message_type}).Warn("Receive group message, but this type is not supported")
	}
}

func groupTextMessage(messageevent *MessageEvent) {
	// get the pure text message, without @xxx
	messageevent.Message.Content = strings.TrimSuffix(strings.TrimPrefix(messageevent.Message.Content, "{\"text\":\""), "\"}")
	messageevent.Message.Content = messageevent.Message.Content[strings.Index(messageevent.Message.Content, " ")+1:]
	logrus.WithFields(logrus.Fields{"message content": messageevent.Message.Content}).Info("Receive group TEXT message")

	if handler, exists := groupMessageMap[messageevent.Message.Content]; exists {
		// "立项"是产品经理群所有人权限
		if(messageevent.Message.Content == "立项"){
			chat_id := 	messageevent.Message.Chat_id
			if(chat_id != LeaderGroupID && chat_id != DevGroupID){
				return
			}
			handler(messageevent)
			return
		} else {
			// 群所有者权限
			chat_id := 	messageevent.Message.Chat_id
			owner ,_ := robot.Robot.GetGroupOwner(messageevent.Message.Chat_id)
			if(chat_id != DevGroupID &&  messageevent.Sender.Sender_id.Open_id != owner){
				return
			}else{
				handler(messageevent)
				return
			}
		}
	} else {
		logrus.Error("Group message failed to find event handler: ", messageevent.Message.Content)
		global.Cli.Send("chat_id", messageevent.Message.Chat_id, "text", "关键词"+" ["+messageevent.Message.Content+"] "+"未定义！")
		return
	}
}

func GroupMessageRegister(f messageHandler, s string) {

	if _, isEventExist := groupMessageMap[s]; isEventExist {
		logrus.Warning("Double declaration of group message handler: ", s)
	}
	groupMessageMap[s] = f
}
