package model

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	_ "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RobotState struct {
	GroupID   string `bson:"group_id"`
	GroupSpace string `bson:"group_space"`
	GroupOwner string `bson:"group_owner"`
}

func insertRobotStateRecord(chatID string, spaceID string, userID string) {
	rs := bson.D{{"group_id", chatID},{"group_space", spaceID},{"group_owner", userID}}

	res , err := robot_state.InsertOne(context.TODO(), rs)
	if err != nil {
		logrus.Info("[robot-state-db] insert robot state record failed with error :", err)
		logrus.Error(err)
	}

	logrus.Info("[robot-state-db] inserted robotStateRecord with ID ", res.InsertedID)
}

func updateRobotStateRecord(chatID string, spaceID string, userID string) {
	
	filter := bson.D{{"group_id", chatID}}
	update := bson.D{{"$set",bson.D{{"group_id", chatID},{"group_space", spaceID},{"group_owner", userID}}}}

	result , err := robot_state.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		logrus.Info("[robot-state-db] update robot state record failed with error :", err)
		logrus.Error(err)
	}

	if result.MatchedCount != 0 {
		logrus.Info("[robot-state-db] matched and replaced an existing document with chatID "+ chatID)
		return
	}
}

func FindRobotStateRecords() (map[string]string, map[string]string) {
	
	cur, err := robot_state.Find(context.TODO(), bson.D{})
	if err != nil {
		logrus.Error(err)
	}

	var results []bson.M
	if err = cur.All(context.TODO(), &results); err != nil {
		logrus.Fatal(err)
	}

	group_space := make(map[string]string)
	group_owner := make(map[string]string)

	for _, result := range results {
		group_space[result["group_id"].(string)] = result["group_space"].(string)
		group_owner[result["group_id"].(string)] = result["group_owner"].(string)
	}

	return group_space, group_owner
}

func DeleteRobotStateRecords(chatID string) {
	
	filter := bson.D{{"group_id", chatID}}

	_, err := robot_state.DeleteMany(context.TODO(), filter)
	if err != nil {
		logrus.Error(err)
	}

	logrus.Info("[robot-state-db] deleted document with chatID :", chatID)
}

func SetGroupSpaceRecord(groupID string, spaceID string) {
	// if groupID is not in the database, insert a new record
	filter := bson.D{{"group_id", groupID}}
	update := bson.D{{"$set", bson.D{{"group_space", spaceID}}}}
	var result bson.M
	err := robot_state.FindOneAndUpdate(context.TODO(), filter, update).Decode(&result)
	if err == mongo.ErrNoDocuments {
		// insert a new record
		insertRobotStateRecord(groupID, spaceID, "")
	}else if err != nil {
		logrus.Info("[robot-state-db] updated group space record failed with error :", err)
		logrus.Error(err)
	}else{
	logrus.Info("[robot-state-db] updated group space record with groupID :", groupID)
	}
}

func SetGroupOwnerRecord(groupID string, userID string) {
	// if groupID is not in the database, insert a new record
	filter := bson.D{{"group_id", groupID}}
	update := bson.D{{"$set", bson.D{{"group_owner", userID}}}}
	var result bson.M
	err := robot_state.FindOneAndUpdate(context.TODO(), filter, update).Decode(&result)
	if err == mongo.ErrNoDocuments {
		// insert a new record
		insertRobotStateRecord(groupID, "", userID)
	}else if err!= nil {
		logrus.Info("[robot-state-db] updated group owner record failed with error :", err)
		logrus.Error(err)
	}else{
		logrus.Info("[robot-state-db] updated group owner record with groupID :", groupID)
	}
}