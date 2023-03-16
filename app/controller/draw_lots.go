package controller

import (
	"encoding/xml"
	"github.com/YasyaKarasu/feishuapi"
	"github.com/sirupsen/logrus"
	"math/rand"
	"strconv"
	"strings"
	"xlab-feishu-robot/global"
	"xlab-feishu-robot/model"
)

type at struct {
	userID   string `json:"user_id"`
	userName string `json:"user_name"`
}

// Define state
const (
	X  = 0 // 未启动
	S0 = 1 // 问人员
	S1 = 2 // 问组数
	S2 = 3 // 问每组人数
)

var currentState = X

// These IDs are all OPEN_ID
var participantsID []string
var mapIDName = make(map[string]string)
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
		err = GetParticipants(messageevent)
		if err != nil {
			InputError(messageevent)
			return
		}
		currentState = S1
		global.Feishu.MessageSend(feishuapi.GroupChatId, groupID, feishuapi.Text, "请问需要抽取几组？")
	case S1:
		count, err = GetNumber(messageevent)
		if err != nil {
			InputError(messageevent)
			return
		}
		currentState = S2
		global.Feishu.MessageSend(feishuapi.GroupChatId, groupID, feishuapi.Text, "请问每组需要几人？")
	case S2:
		size, err := GetNumber(messageevent)
		if err != nil {
			InputError(messageevent)
			return
		}
		// start to draw lots
		groups, err := DrawLots(participantsID, count, size, groupID)
		if err != nil {
			InputError(messageevent)
			return
		}
		currentState = X
		SendResult(groups, groupID)
	}
}

func InputError(messageevent *model.MessageEvent) {
	groupID := messageevent.Message.Chat_id
	global.Feishu.MessageSend(feishuapi.GroupChatId, groupID, feishuapi.Text, "输入格式有误，请重新输入")
	currentState = X
}

// GetParticipants is a function to get participants' ID and name
// The format of message content is like this:
// <at user_id="xxx">xxx</at><at user_id="xxx">xxx</at><at user_id="xxx">xxx</at>
// So we need to use xml.Unmarshal to get the participants' ID and name
func GetParticipants(messageevent *model.MessageEvent) (err error) {
	var ats []at
	err = xml.Unmarshal([]byte(messageevent.Message.Content), &ats)
	if err != nil {
		logrus.Error(err)
		return
	}
	for _, at := range ats {
		participantsID = append(participantsID, at.userID)
		mapIDName[at.userID] = at.userName
	}
	return
}

func GetNumber(messageevent *model.MessageEvent) (number int, err error) {
	number, err = strconv.Atoi(messageevent.Message.Content)
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

// DrawLots is a function to pick [count] groups from participantsID, each group has [size] people
func DrawLots(participantsID []string, count int, size int, groupID string) (groups [][]string, err error) {
	// Check if the number of participants is enough
	if len(participantsID) < count*size {
		global.Feishu.MessageSend(feishuapi.GroupChatId, groupID, feishuapi.Text, "参与人数不足")
		return
	}
	// Pick [count] groups from participantsID randomly
	for i := 0; i < count; i++ {
		var group []string
		for j := 0; j < size; j++ {
			// Pick a random number
			random := rand.Intn(len(participantsID))
			// Pick a person from participantsID randomly
			group = append(group, participantsID[random])
			// Remove the person from participantsID
			participantsID = append(participantsID[:random], participantsID[random+1:]...)
		}
		groups = append(groups, group)
	}
	return
}

func SendResult(groups [][]string, groupID string) {
	// string builder
	var sb strings.Builder
	for i, group := range groups {
		sb.WriteString("第" + strconv.Itoa(i+1) + "组：")
		// @ user in the format of <at user_id="xxx">xxx</at>
		for _, userID := range group {
			sb.WriteString("<at user_id=\"" + userID + "\">" + mapIDName[userID] + "</at>")
		}
		sb.WriteString("\n")
	}
	global.Feishu.MessageSend(feishuapi.GroupChatId, groupID, feishuapi.Text, sb.String())
}
