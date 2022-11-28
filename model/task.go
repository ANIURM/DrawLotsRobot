package model

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	task *mongo.Collection
)

type TaskStatus string

const (
	NotStarted TaskStatus = "未开始"
	InProgress TaskStatus = "进行中"
	Completed  TaskStatus = "已完成"
)

type Task struct {
	TaskId       primitive.ObjectID `bson:"_id,omitempty"`
	TaskRecordId string             `bson:"task_record_id"`

	ProjectChat      string     `bson:"project_chat"`       // chat_id
	TaskName         string     `bson:"task_name"`          // 任务名称
	TaskStatus       TaskStatus `bson:"task_status"`        // 任务状态
	TaskManagerIds   []string   `bson:"task_manager_ids"`   //任务负责人
	TaskManagerNames []string   `bson:"task_manager_names"` //任务负责人
	TaskStartTime    string     `bson:"task_start_time"`    //开始时间
	TaskEndTime      string     `bson:"task_end_time"`      //结束时间
}

func InsertTaskRecords(v []Task) {
	var l []interface{}

	for _, value := range v {
		l = append(l, value)
	}

	_, err := task.InsertMany(context.TODO(), l)
	if err != nil {
		logrus.Error(err)
	}
}

func QueryTaskRecordsByChat(chat_id string) (*[]Task, error) {
	var tasks []Task

	filter := bson.D{{Key: "project_chat", Value: chat_id}}

	cur, err := task.Find(context.TODO(), filter)
	if err != nil {
		logrus.New().WithField("ProjectChat", chat_id).Error(err)
		return nil, err
	}

	var result Task
	for cur.Next(context.TODO()) {
		cur.Decode(&result)
		tasks = append(tasks, result)
	}
	return &tasks, nil
}

func UpdateTaskRecord(ProjectChat string, taskStatus, taskStartTime, taskEndTime string) {

	filter := bson.D{{Key: "project_chat", Value: ProjectChat}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "project_chat", Value: ProjectChat}, {Key: "task_status", Value: taskStatus}, {Key: "task_start_time", Value: taskStartTime}, {Key: "task_end_time", Value: taskEndTime}}}}

	result, err := task.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		logrus.Info("[robot-state-db] update task record failed with error :", err)
		logrus.Error(err)
	}

	if result.MatchedCount != 0 {
		logrus.Info("[robot-state-db] matched and replaced an existing document with chatID " + ProjectChat)
		return
	}
}
