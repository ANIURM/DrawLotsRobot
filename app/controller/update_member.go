package controller

import (
	"strings"
	"xlab-feishu-robot/global"
	"xlab-feishu-robot/model"

	"github.com/YasyaKarasu/feishuapi"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

func ReceiveUpdateMember(msgEvent *model.MessageEvent){
	UpdateMember(msgEvent.Message.Chat_id)
	space_id, err := model.QueryKnowledgeSpaceByChat(msgEvent.Message.Chat_id)
	if err != nil {
		return
	}
	//提醒项目经理填表
	SendProjectManageUrl(msgEvent.Message.Chat_id, space_id)
}

func UpdateMember(group_id string) {
	space_id, err := model.QueryKnowledgeSpaceByChat(group_id)
	if err != nil {
		return
	}

	allEmployees := global.Feishu.EmployeeGetAllInfo(feishuapi.OpenId)
	employeesByName := make(map[string][]string)
	for _, employee := range allEmployees {
		employeesByName[employee.Name] = append(employeesByName[employee.Name], employee.Id)
	}

	logrus.Info(employeesByName)

	allNode := global.Feishu.KnowledgeSpaceGetAllNodes(space_id)
	allNameInRecords := make(map[string]bool)
	for _, node := range allNode {
		if node.Title == "核心成员与职务" {
			allBitables := global.Feishu.DocumentGetAllBitables(node.ObjToken)
			allTables := global.Feishu.DocumentGetAllTables(allBitables[0].AppToken)
			allRecords := global.Feishu.DocumentGetAllRecords(allBitables[0].AppToken, allTables[0].TableId)
			for _, value := range allRecords {
				data := value.Fields
				var name string
				var ok bool
				if name, ok = data["姓名"].(string); !ok {
					logrus.WithField("wiki_url", Url.UrlHead+node.NodeToken).
						Warning("name incomplete")
					continue
				}

				name = strings.Trim(name, "@")
				name = strings.Trim(name, " ")

				allNameInRecords[name] = true
			}
		}
	}

	allNameInGroupList := global.Feishu.GroupGetMembers(group_id,feishuapi.OpenId)
	allNameInGroup := make(map[string]bool)
	for _, memberInfo := range allNameInGroupList{
		allNameInGroup[memberInfo.Name] = true
	}

	//对于 group 中不存在，而 Records 中存在的，要添加成员
	for nameInRecord := range allNameInRecords{
		_, exist := allNameInGroup[nameInRecord]
		if(!exist){
			employeeIds := employeesByName[nameInRecord]
			logrus.Info(employeeIds) 
			global.Feishu.GroupAddMembers(group_id, feishuapi.OpenId, "2", employeeIds)
		}
	}

	//对于 group 中存在，而 Records 中不存在的，要删除成员
	for nameInRecord := range allNameInGroup{
		_, exist := allNameInRecords[nameInRecord]
		if(!exist){
			employeeIds := employeesByName[nameInRecord]
			global.Feishu.GroupDeleteMembers(group_id, feishuapi.OpenId, employeeIds)
		}
	}
}

func SendProjectManageUrl(chatId string, spaceId string) {
	msg := "请项目经理填写甘特图、排期表、任务进度管理：\n"
	var titles []string
	titles = append(titles, "排期甘特图", "项目会议", "任务进度管理")

	nodes := global.Feishu.KnowledgeSpaceGetAllNodes(spaceId)
	for _, value := range nodes {
		if in(value.Title, titles) {
			msg = msg + Url.UrlHead + value.NodeToken + " \n"
		}
		if value.HasChild {
			n := global.Feishu.KnowledgeSpaceGetAllNodes(spaceId, value.NodeToken)
			for _, v := range n {
				if in(v.Title, titles) {
					msg = msg + Url.UrlHead + v.NodeToken + "\n"
				}
			}
		}
	}

	global.Feishu.MessageSend(feishuapi.GroupChatId, chatId, feishuapi.Text, msg)
}

func StartMemberUpdateTimer(groupID string, c *cron.Cron) bool{
	// every day at 0:00
	_, err := c.AddFunc("0 0 0 1/1 * *", func() {
		UpdateMember(groupID)
	})

	if err != nil {
		logrus.Error("[timer] ", groupID, " add member update timer fail")
		logrus.Error(err)
		return true
	}

	return false
}