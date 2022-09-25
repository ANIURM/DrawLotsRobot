package global

import(
	_ "github.com/robfig/cron/v3"
)

// finite state machine
type robotState struct {
	GroupSpace map[string]string
}

var Rob = robotState{
	GroupSpace: make(map[string]string),
}