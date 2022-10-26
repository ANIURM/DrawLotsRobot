package controller

import(
	"xlab-feishu-robot/app/chat"
	"xlab-feishu-robot/global"
	"github.com/sirupsen/logrus"
	"github.com/robfig/cron/v3"
)

var oldRecordInfo []RecordInfo

func ProjectScheduleReminder(messageevent *chat.MessageEvent){
	groupID := messageevent.Message.Chat_id
	checkScheduleUpdated(groupID)
}

func checkScheduleUpdated(groupID string){
	space_id := global.Rob.GroupSpace[groupID]
	_,fileToken := getNodeFileToken(space_id, "排期甘特图", "任务进度管理")
	allBitables := global.Cli.GetAllBitables(fileToken)

	var tableInfoList []TableInfo
	tableInfoList = GetAllTableInfo(allBitables)
	logrus.Debug("[schedule] tableInfoList: ", tableInfoList)

	var recordInfoList []RecordInfo
	recordInfoList = GetAllRecordInfo(tableInfoList)
	logrus.Debug("[shchedule] recordInfoList: ", recordInfoList)

	user_id := global.Rob.GroupOwner[groupID]
	modified := CheckRecordInfoModified(recordInfoList, oldRecordInfo)
	if(modified){
		global.Cli.Send("user_id", user_id, "text", "排期甘特图有更新，请及时查看")
	}else{
		global.Cli.Send("user_id", user_id, "text", "排期甘特图无更新，记得及时修改")
	}
	oldRecordInfo = recordInfoList

}

func getNodeFileToken(space_id string, topFile string, secFile string) (string,string){
	var nodeToken, fileToken string
	allNode := global.Cli.GetAllNodes(space_id)
	for _, node := range allNode {
		if node.Title == topFile {
			allSubNode := global.Cli.GetAllNodes(space_id, node.NodeToken)
			for _, subNode := range allSubNode {
				if subNode.Title == secFile {
					nodeToken = subNode.NodeToken
					fileToken = subNode.ObjToken
					break
				}
			}
			break
		}
	}
	return nodeToken,fileToken
}

func CheckRecordInfoModified(newRecordInfoList []RecordInfo, oldRecordInfoList []RecordInfo) bool{
	if(len(newRecordInfoList) != len(oldRecordInfoList)){
		return true
	}
	for i := 0; i < len(newRecordInfoList); i++{
		if(newRecordInfoList[i].Last_modified_time != oldRecordInfoList[i].Last_modified_time){
			return true
		}
	}
	return false
}

func StartProjectScheduleTimer(groupID string, c *cron.Cron){
	logrus.Info("[timer] add project schedule timer")

	// every two days
	c.AddFunc("0 0 0 */2 * *", func() {
		checkScheduleUpdated(groupID)
	})

}