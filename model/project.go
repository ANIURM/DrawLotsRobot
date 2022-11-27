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

type ProStatus string

const (
	BeforeStart ProStatus = "beforeStart"
	Pending     ProStatus = "pending"
	Revising    ProStatus = "revising"
	Hang        ProStatus = "hang"
	Finished    ProStatus = "finished"
	Aborted     ProStatus = "aborted"
)

type Project struct {
	ProjectId       primitive.ObjectID `bson:"_id,omitempty"`
	ProjectName     string             `bson:"ProjectName"`
	ProjectType     ProjectType        `bson:"ProjectType"`
	ProjectLeaderId string             `bson:"ProjectLeaderId"`
	ProjectSpace    string             `bson:"ProjectSpace"`
	ProjectChat     string             `bson:"ProjectChat"`
	ProjectStatus   ProStatus          `bson:"ProjectStatus,omitempty"`
}

func InsertProjectRecords(v []Project) {
	var l []interface{}
	for _, value := range v {
		l = append(l, value)
	}
	_, err := project.InsertMany(context.TODO(), l)
	if err != nil {
		logrus.Error(err)
	}
}

func QueryProjectRecordsByName(name string) ([]Project, error) {
	var projects []Project

	filter := bson.D{{Key: "ProjectName", Value: name}}

	cur, err := project.Find(context.TODO(), filter)
	if err != nil {
		logrus.WithField("ProjectName", name).Error(err)
		return nil, err
	}

	var result Project
	for cur.Next(context.TODO()) {
		err = cur.Decode(&result)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		projects = append(projects, result)
	}
	return projects, nil
}

func QueryProjectRecordsByType(projtype ProjectType) ([]Project, error) {
	var projects []Project

	filter := bson.D{{Key: "ProjectType", Value: projtype}}

	cur, err := project.Find(context.TODO(), filter)
	if err != nil {
		logrus.WithField("ProjectType", projtype).Error(err)
		return nil, err
	}

	var result Project
	for cur.Next(context.TODO()) {
		err = cur.Decode(&result)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		projects = append(projects, result)
	}
	return projects, nil
}

func QueryProjectRecordsByStatus(status ProStatus) ([]Project, error) {
	var projects []Project

	filter := bson.D{{Key: "ProjectStatus", Value: status}}

	cur, err := project.Find(context.TODO(), filter)
	if err != nil {
		logrus.WithField("ProjectStatus", status).Error(err)
		return nil, err
	}

	var result Project
	for cur.Next(context.TODO()) {
		err = cur.Decode(&result)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		projects = append(projects, result)
	}
	return projects, nil
}

// one chat one project
func QueryProjectRecordsByChat(name string) (Project, error) {
	filter := bson.D{{Key: "ProjectChat", Value: name}}

	cur := project.FindOne(context.TODO(), filter)

	var result Project
	err := cur.Decode(&result)
	if err != nil {
		logrus.Error(err)
		return Project{}, err
	}
	return result, nil
}

func UpdateProjectStatusByChat(pro Project) error {
	filter := bson.D{{Key: "ProjectChat", Value: pro.ProjectChat}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "ProjectStatus", Value: pro.ProjectStatus}}}}

	_, err := project.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		logrus.WithField("ProjectChat", pro.ProjectChat).Error(err)
		return err
	}

	return nil
}

func QueryProjectNameByChat(chatID string) (string, error) {
	filter := bson.D{{Key: "ProjectChat", Value: chatID}}

	cur := project.FindOne(context.TODO(), filter)

	var result Project
	err := cur.Decode(&result)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	return result.ProjectName, nil
}

func QueryProjectLeaderByChat(chatID string) (string, error) {
	filter := bson.D{{Key: "ProjectChat", Value: chatID}}

	cur := project.FindOne(context.TODO(), filter)

	var result Project
	err := cur.Decode(&result)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	return result.ProjectLeaderId, nil
}

func QueryKnowledgeSpaceByChat(chatID string) (string, error) {
	filter := bson.D{{Key: "ProjectChat", Value: chatID}}

	cur := project.FindOne(context.TODO(), filter)

	var result Project
	err := cur.Decode(&result)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	return result.ProjectSpace, nil
}

func QueryChatStatusMap() (map[string]ProStatus, error) {
	cur, err := project.Find(context.TODO(), bson.D{{}})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	var result Project
	var chatStatusMap = make(map[string]ProStatus)
	for cur.Next(context.TODO()) {
		err = cur.Decode(&result)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		chatStatusMap[result.ProjectChat] = result.ProjectStatus
	}

	return chatStatusMap, nil
}
