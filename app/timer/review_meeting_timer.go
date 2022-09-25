package timer

import (
	_ "github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	_ "github.com/YasyaKarasu/feishuapi"
	_ "time"
	"xlab-feishu-robot/pkg/global"
)

func ReviewMeetingTimer() {
	logrus.Info("review meeting timer start")

	//todo: space id
	space_id := "7145117180906979330"
	allNode := global.Cli.GetAllNodes(space_id)
	topFile := "排期甘特图"
	secFile := "任务进度管理"
	var fileToken string

	for _, node := range allNode {
		if node.Title == topFile {
			allSubNode := global.Cli.GetAllNodes(space_id, node.NodeToken)
			for _, subNode := range allSubNode {
				if subNode.Title == secFile {
					// find the corresponding file
					fileToken = subNode.ObjToken
					break
				}
			}
			break
		}
	}

	logrus.Info("file token: ", fileToken)
	logrus.Info("checking: ", fileToken)
	allBlock := global.Cli.GetAllBitables(fileToken)
	logrus.Info("all block: ", allBlock)
	allTable := global.Cli.GetAllTables(fileToken)
	logrus.Info("all table: ", allTable)
	logrus.Info("review meeting timer end")
	return 

}

					