package robot

import (
	"xlab-feishu-robot/model"

	"github.com/sirupsen/logrus"
)

func InitRobotState() {
	Robot.groupSpace, Robot.groupOwner = model.FindRobotStateRecords()
	logrus.Info("[robot] robot state init success")
}
