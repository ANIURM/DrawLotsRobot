package rob

import(
	_ "github.com/robfig/cron/v3"
	"xlab-feishu-robot/model"
)

type robotState struct {
	groupSpace map[string]string // groupID ---> spaceID
	groupOwner map[string]string // groupID ---> userID
}

func (r *robotState) SetGroupSpace(groupID string, spaceID string) {
	r.groupSpace[groupID] = spaceID
	model.SetGroupSpaceRecord(groupID, spaceID)
}

func (r *robotState) GetGroupSpace(groupID string) string {
	return r.groupSpace[groupID]
}

func (r *robotState) SetGroupOwner(groupID string, userID string) {
	r.groupOwner[groupID] = userID
	model.SetGroupOwnerRecord(groupID, userID)
}

func (r *robotState) GetGroupOwner(groupID string) string {
	return r.groupOwner[groupID]
}

func (r *robotState) GetGroupSpaceMap() map[string]string {
	return r.groupSpace
}

var Rob = robotState{
	groupSpace: make(map[string]string),
	groupOwner: make(map[string]string),
}