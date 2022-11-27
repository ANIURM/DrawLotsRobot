package controller

import (
	"strings"
	"xlab-feishu-robot/app/chat"
	"xlab-feishu-robot/global"
	"xlab-feishu-robot/model"

	"github.com/YasyaKarasu/feishuapi"
	"github.com/sirupsen/logrus"
)

func UpdateMember(msgEvent *chat.MessageEvent) {
	space_id, err := model.QueryKnowledgeSpaceByChat(msgEvent.Message.Chat_id)
	if err != nil {
		return
	}

	allEmployees := global.Feishu.GetAllEmployees(feishuapi.OpenId)
	allNode := global.Feishu.GetAllNodes(space_id)
	for _, node := range allNode {
		if node.Title == "核心成员与职务" {
			allBitables := global.Feishu.GetAllBitables(node.ObjToken)
			allTables := global.Feishu.GetAllTables(allBitables[0].AppToken)
			allRecords := global.Feishu.GetAllRecords(allBitables[0].AppToken, allTables[0].TableId)
			for _, value := range allRecords {
				data := value.Fields
				var name string
				var ok bool
				if name, ok = data["姓名"].(string); !ok {
					logrus.WithField("wiki_url", Url.UrlHead+node.NodeToken).
						Warning("name incomplete")
					continue
				}
				if _, ok = data["职务"].(string); !ok {
					logrus.WithField("wiki_url", Url.UrlHead+node.NodeToken).
						Warning("job incomplete")
					continue
				}

				name = strings.Trim(name, "@")
				name = strings.Trim(name, " ")
				employee_id := make([]string, 0)
				for _, v := range allEmployees {
					logrus.Info(v.Name, " ", name, " ", v.Name == name)
					if v.Name == name {
						logrus.Info(v.Id)
						employee_id = append(employee_id, v.Id)
					}
				}
				logrus.Info(employee_id)
				global.Feishu.AddMembers(msgEvent.Message.Chat_id, feishuapi.OpenId, "2", employee_id)
			}
		}
	}

	//提醒项目经理填表
	SendProjectManageUrl(msgEvent.Message.Chat_id, space_id)
}

func SendProjectManageUrl(chatId string, spaceId string) {
	msg := "请项目经理填写甘特图、排期表、任务进度管理：\n"
	var titles []string
	titles = append(titles, "排期甘特图", "项目会议", "任务进度管理")

	nodes := global.Feishu.GetAllNodes(spaceId)
	for _, value := range nodes {
		if in(value.Title, titles) {
			msg = msg + Url.UrlHead + value.NodeToken + " \n"
		}
		if value.HasChild {
			n := global.Feishu.GetAllNodes(spaceId, value.NodeToken)
			for _, v := range n {
				if in(v.Title, titles) {
					msg = msg + Url.UrlHead + v.NodeToken + "\n"
				}
			}
		}
	}

	global.Feishu.Send(feishuapi.GroupChatId, chatId, feishuapi.Text, msg)
}
