package controller

import (
	"encoding/json"
	"fmt"
	"github.com/YasyaKarasu/feishuapi"
	"math/rand"
	"strconv"
	"strings"
	"xlab-feishu-robot/global"
	"xlab-feishu-robot/model"
)

type AtTag struct {
	Tag      string `json:"tag"`
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
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
		groups, err := DrawLots(participantsID, count, size)
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

func GetParticipants(messageevent *model.MessageEvent) (err error) {
	// Unmarshal the JSON input into a slice of AtTag structs
	var atTags []AtTag
	err = json.Unmarshal([]byte(messageevent.Message.Content), &atTags)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Extract the user_id fields for each AtTag object
	for _, atTag := range atTags {
		if atTag.Tag == "at" {
			participantsID = append(participantsID, atTag.UserID)
		}
	}
	return
}

func GetNumber(messageevent *model.MessageEvent) (number int, err error) {
	// Remove the prefix and suffix of the message
	messageevent.Message.Content = strings.TrimSuffix(strings.TrimPrefix(messageevent.Message.Content, "{\"text\":\""), "\"}")
	// Remove spaces in front of the message
	numberString := messageevent.Message.Content[strings.Index(messageevent.Message.Content, " ")+1:]
	// Convert string to int
	number, err = strconv.Atoi(numberString)
	return
}

// DrawLots is a function to pick [count] groups from participantsID, each group has [size] people
func DrawLots(participantsID []string, count int, size int) (groups [][]string, err error) {
	// Check if the number of participants is enough
	if len(participantsID) < count*size {
		err = fmt.Errorf("the number of participants is not enough")
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

}
