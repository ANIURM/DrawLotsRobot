package global

import (
	"github.com/robfig/cron/v3"
)

type GroupTimers struct {
	GTimers map[string] *cron.Cron
}

var Timer = GroupTimers{
	GTimers: make(map[string] *cron.Cron),
}