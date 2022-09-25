package model

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProjectType string

const (
	Internal ProjectType = "internal"
	External ProjectType = "external"
)

type ProjectStatus string

const (
	BeforeStart ProjectStatus = "beforeStart"
	Pending     ProjectStatus = "pending"
	Revising    ProjectStatus = "revising"
	Hang        ProjectStatus = "hang"
	Finished    ProjectStatus = "finished"
	Aborted     ProjectStatus = "aborted"
)

type Project struct {
	ProjectId primitive.ObjectID `bson:"_id,omitempty"`

	ProjectName      string      `bson:"project_name"`
	ProjectType      ProjectType `bson:"project_type"`
	ProjectLeaderIds []string    `bson:"project_leader_ids"`

	GanttDocUrl   string `bson:"gantt_doc_url"`
	PrdDocUrl     string `bson:"prd_doc_url"`
	TechDocUrl    string `bson:"tech_doc_url"`
	FeishuRepoUrl string `bson:"feishu_repo_url"`

	ProjectStatus ProjectStatus `bson:"project_status,omitempty"`
}

func InsertProjectRecords(v []Project) {
	var l []interface{}

	for _, value := range v {
		l = append(l, value)
	}

	project.InsertMany(context.TODO(), l)
}

func QueryProjectRecordsByName(name string) (*[]Project, error) {
	var projects []Project

	filter := bson.D{{Key: "project_name", Value: name}}

	cur, err := project.Find(context.TODO(), filter)
	if err != nil {
		logrus.New().WithField("ProjectName", name).Error(err)
		return nil, err
	}

	var result Project
	for cur.Next(context.TODO()) {
		cur.Decode(&result)
		projects = append(projects, result)
	}
	return &projects, nil
}

func QueryProjectRecordsByType(projtype ProjectType) (*[]Project, error) {
	var projects []Project

	filter := bson.D{{Key: "project_type", Value: projtype}}

	cur, err := project.Find(context.TODO(), filter)
	if err != nil {
		logrus.New().WithField("ProjectType", projtype).Error(err)
		return nil, err
	}

	var result Project
	for cur.Next(context.TODO()) {
		cur.Decode(&result)
		projects = append(projects, result)
	}
	return &projects, nil
}

func QueryProjectRecordsByStatus(status ProjectStatus) (*[]Project, error) {
	var projects []Project

	filter := bson.D{{Key: "project_status", Value: status}}

	cur, err := project.Find(context.TODO(), filter)
	if err != nil {
		logrus.New().WithField("ProjectStatus", status).Error(err)
		return nil, err
	}

	var result Project
	for cur.Next(context.TODO()) {
		cur.Decode(&result)
		projects = append(projects, result)
	}
	return &projects, nil
}
