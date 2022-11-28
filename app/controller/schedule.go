package controller

import (
	"strconv"
	"time"
	"xlab-feishu-robot/global"
	"xlab-feishu-robot/model"

	"github.com/YasyaKarasu/feishuapi"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

var oldRecordInfo []feishuapi.RecordInfo

func ProjectScheduleReminder(messageevent *model.MessageEvent) {
	groupID := messageevent.Message.Chat_id
	checkScheduleUpdated(groupID)
}

func checkScheduleUpdated(groupID string) {
	space_id, err := model.QueryKnowledgeSpaceByChat(groupID)
	if err != nil {
		logrus.Warn("[schedule] ", groupID, " get space id fail")
	}
	_, fileToken := getNodeFileToken(space_id, "排期甘特图", "任务进度管理")
	allBitables := global.Feishu.GetAllBitables(fileToken)

	tableInfoList := GetAllTableInfo(allBitables)
	logrus.Debug("[schedule] tableInfoList: ", tableInfoList)

	recordInfoList := GetAllRecordInfo(tableInfoList)
	logrus.Debug("[shchedule] recordInfoList: ", recordInfoList)

	tasklist, err := model.QueryTaskRecordsByChat(groupID)
	if err != nil {
		logrus.WithField("Group Id", groupID).Error("Query task records by ChatId" + groupID + "failed")
	}
	//播报进度数据获取
	var notStarted []string
	var inProgress []string
	var completed []string
	var updatedTasks []model.Task

	for _, bittable := range allBitables {
		tables := global.Feishu.GetAllTablesInBitable(bittable.AppToken)
		for _, table := range tables {
			records := global.Feishu.GetAllRecordsInTable(bittable.AppToken, table.TableId)
			for _, record := range records {
				var task_name string = ""
				var task_manager string = ""
				var task_manager_ids []string
				var task_manager_names []string
				var task_start_time string = ""
				var task_end_time string = ""
				var task_status string = ""
				if record.Fields["任务状态"] != nil {
					task_status = record.Fields["任务状态"].(string)
				}
				if record.Fields["任务名"] != nil {
					task_name = record.Fields["任务名"].(string)
				}
				if record.Fields["开始时间"] != nil {
					timestamp := int64(record.Fields["开始时间"].(float64))
					start_time := time.Unix(timestamp/1000, 0)
					task_start_time = start_time.Format("2006-01-02 03:04:05")
				}
				if record.Fields["结束时间"] != nil {
					timestamp := int64(record.Fields["结束时间"].(float64))
					end_time := time.Unix(timestamp/1000, 0)
					task_end_time = end_time.Format("2006-01-02 03:04:05")
				}
				if record.Fields["负责人"] != nil {
					temp := record.Fields["负责人"].([]interface{})
					for _, v := range temp {
						temp1 := v.(map[string]interface{})
						task_manager = task_manager + temp1["name"].(string) + " "
						task_manager_ids = append(task_manager_ids, temp1["id"].(string))
						task_manager_names = append(task_manager_names, temp1["name"].(string))
					}
				}
				//if record.Fields["任务状态"] == "未开始" {
				//	notStarted = append(notStarted, record.Fields["任务名"].(string))
				//} else if record.Fields["任务状态"] == "进行中" {
				//	inProgress = append(inProgress, record.Fields["任务名"].(string))
				//} else if record.Fields["任务状态"] == "已完成" {
				//	completed = append(completed, record.Fields["任务名"].(string))
				//}
				if task_status == "未开始" {
					notStarted = append(notStarted, task_name)
				} else if task_status == "进行中" {
					inProgress = append(inProgress, task_name)
				} else if task_status == "已完成" {
					completed = append(completed, task_name)
				}
				//db
				var aNewRecord bool = false
				for _, t := range *tasklist {
					if record.RecordId == t.TaskRecordId {
						aNewRecord = true
					}
				}
				//如果是新记录，插入数据库
				if aNewRecord {
					var task model.Task
					task.ProjectChat = groupID
					task.TaskName = task_name
					task.TaskManagerIds = task_manager_ids
					task.TaskManagerNames = task_manager_names
					task.TaskStartTime = task_start_time
					task.TaskEndTime = task_end_time
					task.TaskRecordId = record.RecordId

					var taskList []model.Task
					taskList = append(taskList, task)
					model.InsertTaskRecords(taskList)
					updatedTasks = append(updatedTasks, task)
					logrus.Info("Task: [ ", task.TaskName, " ] has been inserted into db")
				} else {
					for _, t := range *tasklist {
						if record.RecordId == t.TaskRecordId {
							if task_name == t.TaskName {
								if task_status == string(t.TaskStatus) && task_start_time == t.TaskStartTime && task_end_time == t.TaskEndTime {

								} else {
									var task model.Task
									task.ProjectChat = groupID
									task.TaskName = task_name
									task.TaskManagerIds = task_manager_ids
									task.TaskManagerNames = task_manager_names
									task.TaskStartTime = task_start_time
									task.TaskEndTime = task_end_time
									task.TaskRecordId = record.RecordId
									updatedTasks = append(updatedTasks, task)

									model.UpdateTaskRecord(groupID, task_status, task_start_time, task_end_time)
									logrus.Info("Task: [ ", task.TaskName, " ] has been inserted into db")
								}
							}
							break
						}
					}

				}

			}

		}
	}

	user_id, err := model.QueryProjectLeaderByChat(groupID)
	if err != nil {
		logrus.Warn("[schedule] ", groupID, " get project leader fail")
	}
	modified := CheckRecordInfoModified(recordInfoList, oldRecordInfo)

	groupName, err := model.QueryProjectNameByChat(groupID)
	if err != nil {
		logrus.Warn("[schedule] ", groupID, " get project name fail")
	}
	//自动播报
	var msg string
	msg = groupName + ":\n 当前【任务进度管理】看板任务总数" + strconv.Itoa(len(notStarted)+len(inProgress)+len(completed)) + "个，未开始任务" + strconv.Itoa(len(notStarted)) + "个，进行中任务" + strconv.Itoa(len(inProgress)) + "个，已完成任务" + strconv.Itoa(len(completed)) + "个。\n"

	var newTasks []string
	newTasks = append(newTasks, notStarted...)
	newTasks = append(newTasks, inProgress...)
	newTasks = append(newTasks, completed...)

	if modified {
		var updates string
		for _, task := range updatedTasks {
			var task_manager_names string
			for _, name := range task.TaskManagerNames {
				task_manager_names = task_manager_names + " " + name
			}
			updates = updates + "\"" + task.TaskName + "\n" + "负责人：" + task_manager_names + "\n" + "时间：" + task.TaskStartTime + " - " + task.TaskEndTime + "\n" + "状态：" + string(task.TaskStatus) + "\"\n"
		}

		msg = msg + "更新内容：\n" + updates

	} else {
		msg = "请及时更新排期甘特图和任务进度管理，跟进项目进程！ \n" + msg
	}
	link := getlink(space_id)
	msg = msg + "欲了解详细内容，请点击: " + link

	global.Feishu.Send(feishuapi.UserOpenId, user_id, feishuapi.Text, msg)

	oldRecordInfo = recordInfoList

}

func getNodeFileToken(space_id string, topFile string, secFile string) (string, string) {
	var nodeToken, fileToken string
	allNode := global.Feishu.GetAllNodes(space_id)
	for _, node := range allNode {
		if node.Title == topFile {
			allSubNode := global.Feishu.GetAllNodes(space_id, node.NodeToken)
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
	return nodeToken, fileToken
}

func CheckRecordInfoModified(newRecordInfoList []feishuapi.RecordInfo, oldRecordInfoList []feishuapi.RecordInfo) bool {
	if len(newRecordInfoList) != len(oldRecordInfoList) {
		return true
	}
	for i := 0; i < len(newRecordInfoList); i++ {
		if newRecordInfoList[i].LastModifiedTime != oldRecordInfoList[i].LastModifiedTime {
			return true
		}
	}
	return false
}

func StartProjectScheduleTimer(groupID string, c *cron.Cron) bool {

	// every two days at 9:00
	_, err := c.AddFunc("* * 9 1/2 * *", func() {
		checkScheduleUpdated(groupID)
	})

	if err != nil {
		logrus.Error("[timer] ", groupID, " add project schedule timer fail")
		logrus.Error(err)
		return true
	}

	return false
}

func getlink(spaceId string) string {
	var msg string
	var titles []string
	titles = append(titles, "排期甘特图", "任务进度管理")

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
	return msg
}
