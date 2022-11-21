package robot

import (
	"xlab-feishu-robot/model"

	_ "github.com/robfig/cron/v3"
)

type robotState struct {
	groupSpace map[string]string // groupID ---> spaceID
	groupOwner map[string]string // groupID ---> openID
}

func (r *robotState) SetGroupSpace(groupID string, spaceID string) {
	r.groupSpace[groupID] = spaceID
	model.SetGroupSpaceRecord(groupID, spaceID)
}

func (r *robotState) GetGroupSpace(groupID string) (string, bool) {
	v, ok := r.groupSpace[groupID]
	return v, ok
}

func (r *robotState) SetGroupOwner(groupID string, userID string) {
	r.groupOwner[groupID] = userID
	model.SetGroupOwnerRecord(groupID, userID)
}

func (r *robotState) GetGroupOwner(groupID string) (string, bool) {
	v, ok := r.groupOwner[groupID]
	return v, ok
}

func (r *robotState) GetGroupSpaceMap() map[string]string {
	return r.groupSpace
}

func (r *robotState) DeleteGroup(groupID string) {
	delete(r.groupSpace, groupID)
	delete(r.groupOwner, groupID)
	model.DeleteRobotStateRecords(groupID)
}

var Robot = robotState{
	groupSpace: make(map[string]string),
	groupOwner: make(map[string]string),
}
