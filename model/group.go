package model

import (
	"context"
	"xlab-feishu-robot/pkg/global"

	"github.com/YasyaKarasu/feishuapi"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Group struct {
	GroupId primitive.ObjectID `bson:"_id,omitempty"`

	FeishuChatId string `bson:"feishu_chatid,omitempty"`

	GroupName    string   `bson:"group_name"`
	GroupMembers []string `bson:"group_members"`
	GroupOwner   string   `bson:"group_owner"`
	GroupAdmins  []string `bson:"group_admins"`

	ProjectId primitive.ObjectID `bson:"project_id"`

	GroupSettingsJSON string `bson:"group_settings_json"`
}

func getGroupInfo(v *Group, userIdType feishuapi.UserIdType, userAccessToken string) {
	query := make(map[string]string)
	query["user_id_type"] = string(userIdType)

	resp := global.Cli.Request("get", "open-apis/im/v1/chats/"+v.FeishuChatId, query, nil, nil)

	v.GroupName = resp["name"].(string)
	v.GroupOwner = resp["owner_id"].(string)

	l := global.Cli.GetGroupMembers(v.FeishuChatId, userIdType)
	for _, value := range l {
		v.GroupMembers = append(v.GroupMembers, value.MemberId)
	}

	logrus.Info(l)

	body := make(map[string][]string)
	body["manager_ids"] = []string{v.GroupOwner}
	memquery := make(map[string]string)
	memquery["member_id_type"] = string(userIdType)
	header := make(map[string]string)
	header["Authorization"] = userAccessToken
	resp = global.Cli.Request("post", "open-apis/im/v1/chats/"+v.FeishuChatId+"/managers/add_managers", memquery, header, body)

	managers := resp["chat_managers"].([]string)
	v.GroupAdmins = append(v.GroupAdmins, managers...)

	managers = resp["chat_bot_managers"].([]string)
	v.GroupAdmins = append(v.GroupAdmins, managers...)
}

func InsertGroupRecords(v []Group, userIdType feishuapi.UserIdType, userAccessToken string) {
	var l []interface{}

	for _, value := range v {
		getGroupInfo(&value, userIdType, userAccessToken)
		l = append(l, value)
	}

	group.InsertMany(context.TODO(), l)
}

func QueryGroupByChatId(chatid string) (*Group, error) {
	var querygroup Group

	filter := bson.D{{Key: "feishu_chatid", Value: chatid}}
	err := group.FindOne(context.TODO(), filter).Decode(&querygroup)
	if err != nil {
		logrus.WithField("feishu_chatid", chatid).Error(err)
	}
	return &querygroup, nil
}

func QueryGroupByName(name string) (*[]Group, error) {
	var groups []Group

	filter := bson.D{{Key: "group_name", Value: name}}

	cur, err := group.Find(context.TODO(), filter)
	if err != nil {
		logrus.New().WithField("GroupName", name).Error(err)
		return nil, err
	}

	var result Group
	for cur.Next(context.TODO()) {
		cur.Decode(&result)
		groups = append(groups, result)
	}
	return &groups, nil
}

func UpdateGroupSettingsJSON(id primitive.ObjectID, JSON string) error {
	filter := bson.D{{Key: "_id", Value: id}}
	replacement := bson.D{{Key: "group_settings_json", Value: JSON}}

	_, err := group.ReplaceOne(context.TODO(), filter, replacement)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"id":                id,
			"GroupSettingsJSON": JSON,
		}).Error(err)
		return err
	}
	return nil
}
