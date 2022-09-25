package model

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PrivilegeGroup struct {
	PrivilegeGroupId primitive.ObjectID `bson:"_id,omitempty"`

	PrivilegeName string `bson:"privilege_name,omitempty"`
	Explanation   string `bson:"explanation,omitempty"`

	ValidUntil primitive.DateTime `bson:"valid_until"`
}

func InsertPrivilegeRecords(v []PrivilegeGroup) {
	var l []interface{}

	for _, value := range v {
		l = append(l, value)
	}

	privilege.InsertMany(context.TODO(), l)
}

func QueryPrivilegeGroupByName(name string) (*[]PrivilegeGroup, error) {
	var privileges []PrivilegeGroup

	filter := bson.D{{Key: "privilege_name", Value: name}}

	cur, err := privilege.Find(context.TODO(), filter)
	if err != nil {
		logrus.New().WithField("PrivilegeName", name).Error(err)
		return nil, err
	}

	var result PrivilegeGroup
	for cur.Next(context.TODO()) {
		cur.Decode(&result)
		privileges = append(privileges, result)
	}
	return &privileges, nil
}

func UpdatePrivilegeGroupName(id primitive.ObjectID, name string) error {
	filter := bson.D{{Key: "_id", Value: id}}
	replacement := bson.D{{Key: "privilege_name", Value: name}}

	_, err := privilege.ReplaceOne(context.TODO(), filter, replacement)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"id":            id,
			"PrivilegeName": name,
		}).Error(err)
		return err
	}
	return nil
}

func UpdatePrivilegeGroupExplanation(id primitive.ObjectID, explanation string) error {
	filter := bson.D{{Key: "_id", Value: id}}
	replacement := bson.D{{Key: "explanation", Value: explanation}}

	_, err := privilege.ReplaceOne(context.TODO(), filter, replacement)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"id":          id,
			"Explanation": explanation,
		}).Error(err)
		return err
	}
	return nil
}
