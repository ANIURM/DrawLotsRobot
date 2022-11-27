package controller

// this is a simple test file

import (
	"strings"
	"time"
	"xlab-feishu-robot/app/chat"
	"xlab-feishu-robot/global"
	"xlab-feishu-robot/model"

	"github.com/YasyaKarasu/feishuapi"
	"github.com/sirupsen/logrus"
)

// suppose we have already know the space id
func FinishProject(messageevent *chat.MessageEvent) {

	space_id, err := model.QueryKnowledgeSpaceByChat(messageevent.Message.Chat_id)
	if err != nil {
		return
	}

	allNode := global.Feishu.GetAllNodes(space_id)
	requirement := map[string]int{"项目介绍": 100, "产品需求文档": 200, "产品测试记录": 200, "用户手册": 300}
	tooShort := []string{}

	for _, node := range allNode {
		if requireSize, exist := requirement[node.Title]; exist {
			logrus.Trace("checking ", node.Title, " required size: ", requireSize)
			size := recursiveCountNodeSize(space_id, &node)
			if size < requireSize {
				tooShort = append(tooShort, node.Title)
			}
			delete(requirement, node.Title)
		}
		if _, exist := requirement[strings.Trim(node.Title, "$")]; exist {
			logrus.Trace("skip ", node.Title)
			delete(requirement, strings.Trim(node.Title, "$"))
		}
	}

	if len(requirement) != 0 {
		text := "以下文档不存在："
		for key, _ := range requirement {
			text += "<" + key + "> "
		}
		text += "请添加相应文档"
		global.Feishu.Send(feishuapi.GroupChatId, messageevent.Message.Chat_id, feishuapi.Text, text)
	}

	if len(tooShort) != 0 {
		text := "以下文档内容过少："
		for _, node := range tooShort {
			text += "<" + node + ">"
		}
		text += "请补充相应内容"
		global.Feishu.Send(feishuapi.GroupChatId, messageevent.Message.Chat_id, feishuapi.Text, text)
	}

	// success
	if len(requirement) == 0 && len(tooShort) == 0 {
		global.Feishu.Send(feishuapi.GroupChatId, messageevent.Message.Chat_id, feishuapi.Text, "结项成功")
		EndGroupTimer(messageevent.Message.Chat_id)
		// change data in db
		project, err := model.QueryProjectRecordsByChat(messageevent.Message.Chat_id)
		if err != nil {
			logrus.Error(err)
		} else {
			project.ProjectStatus = model.Finished
			model.UpdateProjectStatusByChat(project)
		}
	}
}

func recursiveCountNodeSize(space_id string, node *feishuapi.NodeInfo) int {
	size := 0
	if node.HasChild {
		allNode := global.Feishu.GetAllNodes(space_id, node.NodeToken)
		for _, node := range allNode {
			size += recursiveCountNodeSize(space_id, &node)
		}
	}

	content := global.Feishu.GetRawContent(node.ObjToken)

	size += len(content) / 3 // 3 bytes per chinese character
	logrus.WithFields(logrus.Fields{"resp": content}).Trace("the ", node.Title, " size is: ", size)
	// api 频率限定为每秒 5 次，所以这里需要 sleep 200ms
	time.Sleep(200 * time.Millisecond)
	return size
}
