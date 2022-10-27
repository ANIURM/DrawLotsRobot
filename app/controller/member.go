package controller

import (
	"strings"
	"xlab-feishu-robot/app/chat"
	"xlab-feishu-robot/global"
	"xlab-feishu-robot/model"

	"github.com/YasyaKarasu/feishuapi"
	"github.com/sirupsen/logrus"
)

func initUpdateMember() {
	chat.GroupMessageRegister(UpdateMember, "人员变更")
}

func UpdateMember(msgEvent *chat.MessageEvent) {
	//space_id := global.Rob.GroupSpace[msgEvent.Message.Chat_id]
	space_id := "7141190444620513282"

	allNode := global.Cli.GetAllNodes(space_id)
	for _, node := range allNode {
		if node.Title == "核心成员与职务" {
			allBitables := global.Cli.GetAllBitables(node.ObjToken)
			allTables := global.Cli.GetAllTables(allBitables[0].AppToken)
			allRecords := global.Cli.GetAllRecords(allBitables[0].AppToken, allTables[0].TableId)
			for _, value := range allRecords {
				data := value.Fields
				var name string
				var ok bool
				if name, ok = data["姓名"].(string); !ok {
					logrus.WithField("wiki_url", "https://xn4zlkzg4p.feishu.cn/wiki/"+node.NodeToken).
						Error("name incomplete")
					continue
				}
				if _, ok = data["职务"].(string); !ok {
					logrus.WithField("wiki_url", "https://xn4zlkzg4p.feishu.cn/wiki/"+node.NodeToken).
						Error("job incomplete")
					continue
				}

				name = strings.Trim(name, "@")
				employee, err := model.QueryEmployeeByFullname(name)
				if err != nil {
					logrus.WithField("name", name).
						Error("Query employee info fail")
				}
				employee_id := make([]string, 0)
				for _, value := range *employee {
					employee_id = append(employee_id, value.FeishuOpenId)
				}
				global.Cli.AddMembers(msgEvent.Message.Chat_id, feishuapi.OpenId, "2", employee_id)
			}
		}
	}
}
