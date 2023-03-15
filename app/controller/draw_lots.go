package controller

import (
	"github.com/YasyaKarasu/feishuapi"
	"xlab-feishu-robot/global"
	"xlab-feishu-robot/model"
)

// Define state
const (
	X  = 0 // 未启动
	S0 = 1 // 问人员
	S1 = 2 // 问组数
	S2 = 3 // 问每组人数
)

var currentState = X
var participantsID []string
var count int

// DrawLotsRobot is a function to draw lots
func DrawLotsRobot(messageevent *model.MessageEvent) {
	groupID := messageevent.Message.Chat_id
	var err error
	switch currentState {
	case X:
		currentState = S0
		global.Feishu.MessageSend(feishuapi.GroupChatId, groupID, feishuapi.Text, "请输入参与人员")
	case S0:
		participantsID, err = GetParticipants(messageevent)
		if err != nil {
			global.Feishu.MessageSend(feishuapi.GroupChatId, groupID, feishuapi.Text, "输入格式有误，请重新输入")
			currentState = X
			return
		}
		currentState = S1
		global.Feishu.MessageSend(feishuapi.GroupChatId, groupID, feishuapi.Text, "请问需要抽取几组？")
	case S1:
		count, err = GetGroupCount(messageevent)
		if err != nil {
			global.Feishu.MessageSend(feishuapi.GroupChatId, groupID, feishuapi.Text, "输入格式有误，请重新输入")
			currentState = X
			return
		}
		currentState = S2
		global.Feishu.MessageSend(feishuapi.GroupChatId, groupID, feishuapi.Text, "请问每组需要几人？")
	case S2:
		size, err := GetGroupSize(messageevent)
		if err != nil {
			global.Feishu.MessageSend(feishuapi.GroupChatId, groupID, feishuapi.Text, "输入格式有误，请重新输入")
			currentState = X
			return
		}
		// start draw lots
		groups, err := DrawLots(participantsID, count, size)
		if err != nil {
			global.Feishu.MessageSend(feishuapi.GroupChatId, groupID, feishuapi.Text, "抽签失败，请重新输入")
			currentState = X
			return
		}
		SendResult(groups, groupID)
	}
}

func GetParticipants(messageevent *model.MessageEvent) (participantsID []string, err error) {
	return
}

func GetGroupCount(messageevent *model.MessageEvent) (count int, err error) {
	return
}

func GetGroupSize(messageevent *model.MessageEvent) (size int, err error) {
	return
}

func DrawLots(participantsID []string, count int, size int) (groups [][]string, err error) {
	return
}

func SendResult(groups [][]string, groupID string) {

}
