package controller

// this is a simple test file

import (
	"strings"
	"time"
	"xlab-feishu-robot/app/chat"
	"xlab-feishu-robot/global"
	"xlab-feishu-robot/model"

	"xlab-feishu-robot/global/robot"

	"github.com/YasyaKarasu/feishuapi"
	"github.com/sirupsen/logrus"
)

// suppose we have already know the space id
func ProjectOver(messageevent *chat.MessageEvent) {
	logrus.Debug("project over")

	space_id, _ := robot.Robot.GetGroupSpace(messageevent.Message.Chat_id)
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
		if _, exist := requirement[strings.Trim(node.Title, "$")]; exist{
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

	// success
	if len(requirement) == 0 && len(tooShort) == 0 {
		global.Cli.Send("chat_id", messageevent.Message.Chat_id, "text", "结项成功")
		EndGroupTimer(messageevent.Message.Chat_id)
		robot.Robot.DeleteGroup(messageevent.Message.Chat_id)
		// change data in db
		projectListPointer,err := model.QueryProjectRecordsByChat(messageevent.Message.Chat_id)
		projectList := *projectListPointer
		if(err!= nil){
			logrus.Error(err)
		}else{
			projectList[0].ProjectStatus = model.Finished
			model.UpdateProjectStatusByChat(projectList[0])
		}
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