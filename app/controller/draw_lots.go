package controller

import (
	"errors"
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

// Using map to store participants
// {groupID: []participant}
var participants map[string][]participant
var count int

// DrawLotsRobot is a function to draw lots
func DrawLotsRobot(messageevent *model.MessageEvent) {
	groupID := messageevent.Message.Chat_id
	// Check whether the robot needs to be reset
	if needReset(messageevent.Message.Content) {
		reset(groupID)
		return
	}
	// Check whether the user need help
	if needHelp(messageevent.Message.Content) {
		help(groupID)
		return
	}

	var err error
	// If groupID is not in stateMap, then stateMap[groupID] will be 0(X)
	switch stateMap[groupID] {
	case X:
		stateMap[groupID] = S0
		global.Feishu.MessageSend(feishuapi.GroupChatId, groupID, feishuapi.Text, "请输入参与人员")
	case S0:
		err = getParticipants(messageevent)
		if err != nil {
			inputError(messageevent, err)
			return
		}
		stateMap[groupID] = S1
		global.Feishu.MessageSend(feishuapi.GroupChatId, groupID, feishuapi.Text, "请问需要抽取几组？")
	case S1:
		count, err = getNumber(messageevent)
		if err != nil {
			inputError(messageevent, err)
			return
		}
		stateMap[groupID] = S2
		global.Feishu.MessageSend(feishuapi.GroupChatId, groupID, feishuapi.Text, "请问每组需要几人？")
	case S2:
		size, err := getNumber(messageevent)
		if err != nil {
			inputError(messageevent, err)
			return
		}
		// start to draw lots
		groups, err := drawLots(participants, count, size, groupID)
		if err != nil {
			inputError(messageevent, err)
			return
		}
		stateMap[groupID] = X
		sendResult(groups, groupID)
	}
}

func inputError(messageevent *model.MessageEvent, err error) {
	groupID := messageevent.Message.Chat_id
	global.Feishu.MessageSend(feishuapi.GroupChatId, groupID, feishuapi.Text, "输入格式有误，请重新输入")
	stateMap[groupID] = X
	logrus.Error(err)
}

// getParticipants is a function to get participants' ID and name
func getParticipants(messageevent *model.MessageEvent) (err error) {
	groupID := messageevent.Message.Chat_id
	// Clear last-time participants
	participants[groupID] = nil
	if wantAllMembers(messageevent.Message.Content) {
		getAllGroupMembers(groupID)
	} else {
		getMentionedPerson(messageevent)
	}
	if len(participants[groupID]) == 0 {
		err = errors.New("no participants")
	}
	return
}

func wantAllMembers(content string) bool {
	// Check if the message content contains "所有人" 或者 "all"
	return strings.Contains(content, "所有人") || strings.Contains(content, "all")
}

func getMentionedPerson(messageevent *model.MessageEvent) {
	groupID := messageevent.Message.Chat_id
	// Get openID and name of participants from messageevent.Message.Mentions
	// Skip the first mention, which is the robot itself
	for _, mention := range messageevent.Message.Mentions[1:] {
		participants[groupID] = append(participants[groupID], participant{mention.Id.Open_id, mention.Name})
	}
}

func getAllGroupMembers(groupID string) {
	allMembers := global.Feishu.GroupGetMembers(groupID, feishuapi.OpenId)
	for _, member := range allMembers {
		participants[groupID] = append(participants[groupID], participant{member.MemberId, member.Name})
	}
}

func getNumber(messageevent *model.MessageEvent) (number int, err error) {
	number, err = strconv.Atoi(messageevent.Message.Content)
	if err != nil {
		logrus.Error(err)
		return
	}
	if number <= 0 {
		err = errors.New("number is less than 0")
	}
	return
}

// drawLots is a function to pick [count] groups from participantsID, each group has [size] people
func drawLots(participants map[string][]participant, count int, size int, groupID string) (groups [][]participant, err error) {
	// TODO: Remove duplicate participants

	// Check if the number of participants is enough
	if len(participants[groupID]) < count*size {
		global.Feishu.MessageSend(feishuapi.GroupChatId, groupID, feishuapi.Text, "参与人数不足")
		return
	}
	// Pick [count] groups from participantsID randomly
	for i := 0; i < count; i++ {
		var group []participant
		for j := 0; j < size; j++ {
			// Pick a random number
			random := rand.Intn(len(participants[groupID]))
			// Pick a person from participantsID randomly
			group = append(group, participants[groupID][random])
			// Remove the person from participants
			participants[groupID] = append(participants[groupID][:random], participants[groupID][random+1:]...)
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

// needReset is a function to check whether the robot needs to be reset
// Return true if the robot needs to be reset, otherwise return false
func needReset(content string) bool {
	if strings.Contains(content, "reset") || strings.Contains(content, "重置") {
		return true
	} else {
		return false
	}
}

func reset(groupID string) {
	stateMap[groupID] = X
}

func needHelp(content string) bool {
	if strings.Contains(content, "help") || strings.Contains(content, "帮助") {
		return true
	} else {
		return false
	}
}

func help(groupID string) {
	helpMessage := "抽签机器人使用说明：\n"
	helpMessage += "1. @机器人，启动抽签\n"
	helpMessage += "2. 输入@机器人 所有人 或者 @机器人 all，抽取所有人。输入@机器人 @xxx @xxx，抽取@的人\n"
	helpMessage += "3. 输入@机器人 组数\n"
	helpMessage += "4. 输入@机器人 每组人数\n"
	helpMessage += "即可获得抽签结果\n"
	helpMessage += "输入@机器人 reset 或者 @机器人 重置，重置抽签机器人\n"
	helpMessage += "输入@机器人 help 或者 @机器人 帮助，查看抽签机器人使用说明\n"
	global.Feishu.MessageSend(feishuapi.GroupChatId, groupID, feishuapi.Text, helpMessage)
}
