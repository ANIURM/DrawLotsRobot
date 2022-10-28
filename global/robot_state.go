package global

import(
	_ "github.com/robfig/cron/v3"
)

type robotState struct {
	GroupSpace map[string]string // groupID ---> spaceID
	GroupOwner map[string]string // groupID ---> userID
}

var Rob = robotState{
	GroupSpace: make(map[string]string),
	GroupOwner: make(map[string]string),
}