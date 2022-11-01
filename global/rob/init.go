package rob

import (
	"xlab-feishu-robot/model"
	"github.com/sirupsen/logrus"
)

func InitRobState( ){
	Rob.groupSpace, Rob.groupOwner = model.FindRobotStateRecords()
	logrus.Info("[rob] robot state init success with ", Rob.groupSpace, Rob.groupOwner)
}