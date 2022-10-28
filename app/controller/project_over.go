package controller

// this is a simple test file

import (
	"time"
	"xlab-feishu-robot/app/chat"
	"xlab-feishu-robot/global"

	"github.com/YasyaKarasu/feishuapi"
	"github.com/sirupsen/logrus"
)

// suppose we have already know the space id
func ProjectOver(messageevent *chat.MessageEvent) {
	logrus.Debug("project over")

	space_id := global.Rob.GroupSpace[messageevent.Message.Chat_id]
	allNode := global.Cli.GetAllNodes(space_id)
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
	}

	if len(requirement) != 0 {
		text := "以下文档不存在："
		for key, _ := range requirement {
			text += "<" + key + "> "
		}
		text += "请添加相应文档"
		global.Cli.Send("chat_id", messageevent.Message.Chat_id, "text", text)
	}

	if len(tooShort) != 0 {
		text := "以下文档内容过少："
		for _, node := range tooShort {
			text += "<" + node + ">"
		}
		text += "请补充相应内容"
		global.Cli.Send("chat_id", messageevent.Message.Chat_id, "text", text)
	}

	if len(requirement) == 0 && len(tooShort) == 0 {
		global.Cli.Send("chat_id", messageevent.Message.Chat_id, "text", "结项成功")
	}
}

func recursiveCountNodeSize(space_id string, node *feishuapi.NodeInfo) int {
	size := 0
	if node.HasChild {
		allNode := global.Cli.GetAllNodes(space_id, node.NodeToken)
		for _, node := range allNode {
			size += recursiveCountNodeSize(space_id, &node)
		}
	}

	methon := "GET"
	path := "/open-apis/docx/v1/documents/" + node.ObjToken + "/raw_content"
	query := map[string]string{"document_id": node.ObjToken}
	headers := map[string]string{}
	body := map[string]string{}
	resp := global.Cli.Request(methon, path, query, headers, body)
	size += len(resp["content"].(string)) / 3 // 3 bytes per chinese character
	logrus.WithFields(logrus.Fields{"resp": resp["content"].(string)}).Trace("the ", node.Title, " size is: ", size)
	// api 频率限定为每秒 5 次，所以这里需要 sleep 200ms
	time.Sleep(200 * time.Millisecond)
	return size
}
