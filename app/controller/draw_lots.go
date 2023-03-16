package controller

import (
	"encoding/json"
	"fmt"
	"github.com/YasyaKarasu/feishuapi"
	"github.com/sirupsen/logrus"
	"math/rand"
	"strconv"
	"strings"
	"xlab-feishu-robot/global"
	"xlab-feishu-robot/model"
)

type atTag struct {
	tag      string `json:"tag"`
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

func GetParticipants(messageevent *model.MessageEvent) (err error) {
	// Unmarshal the JSON input into a slice of atTag structs
	var atTags []atTag
	err = json.Unmarshal([]byte(messageevent.Message.Content), &atTags)
	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err}).Error("Unmarshal JSON input")
		return
	}

	// Extract the user_id fields for each atTag object
	for _, atTag := range atTags {
		if atTag.tag == "at" {
			participantsID = append(participantsID, atTag.userID)
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
	query := make(map[string]string)
	query["receive_id_type"] = string(feishuapi.GroupChatId)

	// create the content array
	content := make([]interface{}, 0)

	// loop through the groups and add group labels and IDs to the content array
	for i, group := range groups {
		// add group label
		content = append(content, map[string]interface{}{
			"tag":  "text",
			"text": fmt.Sprintf("group%d: ", i+1),
		})

		// loop through the group IDs and add them to the content array
		for _, id := range group {
			content = append(content, map[string]interface{}{
				"tag":     "at",
				"user_id": id,
			})
		}
	}

	// create the final payload
	payload := map[string]interface{}{
		"post": map[string]interface{}{
			"zh_cn": map[string]interface{}{
				"content": content,
			},
		},
	}
	// encode the payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		logrus.Error(err)
	}

	body := make(map[string]string)
	body["receive_id"] = groupID
	body["content"] = string(jsonData)
	body["msg_type"] = "post"

	resp := feishuapi.AppClient{}.Request("post", "open-apis/im/v1/messages", query, nil, body)
	if resp == nil {
		logrus.WithFields(logrus.Fields{
			"ReceiveID": groupID,
		}).Error("Send post failed")
	}
}
