package controller

import (
	"github.com/YasyaKarasu/feishuapi"
	"github.com/sirupsen/logrus"
	"math/rand"
	"strconv"
	"strings"
	"xlab-feishu-robot/global"
	"xlab-feishu-robot/model"
)

type participant struct {
	openID string
	name   string
}

// Define state
const (
	X  = 0 // 未启动
	S0 = 1 // 问人员
	S1 = 2 // 问组数
	S2 = 3 // 问每组人数
)

// Using map to store state
// {groupID: state}
var stateMap = map[string]int{}

var participants []participant
var count int

// DrawLotsRobot is a function to draw lots
func DrawLotsRobot(messageevent *model.MessageEvent) {
	groupID := messageevent.Message.Chat_id
	var err error
	// If groupID is not in stateMap, then stateMap[groupID] will be 0(X)
	switch stateMap[groupID] {
	case X:
		stateMap[groupID] = S0
		global.Feishu.MessageSend(feishuapi.GroupChatId, groupID, feishuapi.Text, "请输入参与人员")
	case S0:
		err = getParticipants(messageevent)
		if err != nil {
			inputError(messageevent)
			return
		}
		stateMap[groupID] = S1
		global.Feishu.MessageSend(feishuapi.GroupChatId, groupID, feishuapi.Text, "请问需要抽取几组？")
	case S1:
		count, err = getNumber(messageevent)
		if err != nil {
			inputError(messageevent)
			return
		}
		stateMap[groupID] = S2
		global.Feishu.MessageSend(feishuapi.GroupChatId, groupID, feishuapi.Text, "请问每组需要几人？")
	case S2:
		size, err := getNumber(messageevent)
		if err != nil {
			inputError(messageevent)
			return
		}
		// start to draw lots
		groups, err := drawLots(participants, count, size, groupID)
		if err != nil {
			inputError(messageevent)
			return
		}
		stateMap[groupID] = X
		sendResult(groups, groupID)
	}
}

func inputError(messageevent *model.MessageEvent) {
	groupID := messageevent.Message.Chat_id
	global.Feishu.MessageSend(feishuapi.GroupChatId, groupID, feishuapi.Text, "输入格式有误，请重新输入")
	stateMap[groupID] = X
}

// getParticipants is a function to get participants' ID and name
func getParticipants(messageevent *model.MessageEvent) (err error) {
	// Clear last-time participants
	participants = nil
	// Get openID and name of participants from messageevent.Message.Mentions
	// Skip the first mention, which is the robot itself
	for _, mention := range messageevent.Message.Mentions[1:] {
		participants = append(participants, participant{mention.Id.Open_id, mention.Name})
	}
	return
}

func getNumber(messageevent *model.MessageEvent) (number int, err error) {
	number, err = strconv.Atoi(messageevent.Message.Content)
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

// drawLots is a function to pick [count] groups from participantsID, each group has [size] people
func drawLots(participants []participant, count int, size int, groupID string) (groups [][]participant, err error) {
	// TODO: Remove duplicate participants

	// Check if the number of participants is enough
	if len(participants) < count*size {
		global.Feishu.MessageSend(feishuapi.GroupChatId, groupID, feishuapi.Text, "参与人数不足")
		return
	}
	// Pick [count] groups from participantsID randomly
	for i := 0; i < count; i++ {
		var group []participant
		for j := 0; j < size; j++ {
			// Pick a random number
			random := rand.Intn(len(participants))
			// Pick a person from participantsID randomly
			group = append(group, participants[random])
			// Remove the person from participants
			participants = append(participants[:random], participants[random+1:]...)
		}
		groups = append(groups, group)
	}
	return
}

func sendResult(groups [][]participant, groupID string) {
	// string builder
	var sb strings.Builder
	for i, group := range groups {
		sb.WriteString("第" + strconv.Itoa(i+1) + "组：")
		// @ user in the format of <at user_id="xxx">xxx</at>
		for _, person := range group {
			sb.WriteString("<at user_id=\"" + person.openID + "\">" + person.name + "</at>")
		}
		sb.WriteString("\n")
	}
	global.Feishu.MessageSend(feishuapi.GroupChatId, groupID, feishuapi.Text, sb.String())
}
