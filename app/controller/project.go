package controller

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	global "xlab-feishu-robot/pkg/global"
)

var (
	P         FeishuProjectFormPath
	MyProject NewProject
)

type FeishuProjectFormPath struct {
	AppToken string
	TableId  string
}

type NewProject struct {
	Code int64  `json:"code"`
	Data Data   `json:"data"`
	Msg  string `json:"msg"`
}

type Data struct {
	Record Record `json:"record"`
}

type Record struct {
	Fields   Fields `json:"fields"`
	ID       string `json:"id"`
	RecordID string `json:"record_id"`
}

type Fields struct {
	ProjectName          string                `json:"项目名称"`
	ProjectProperties    string                `json:"项目属性"`
	ProjectSource        string                `json:"项目来源"`
	ProjectProfile       string                `json:"项目简介"`
	ParticipatingMembers []ParticipatingMember `json:"主要参与人员"`
	ProjectManager       []ParticipatingMember `json:"产品经理"`
	//CreatTime            string                `json:"创建时间"`
}

type ParticipatingMember struct {
	Email  string `json:"email"`
	EnName string `json:"en_name"`
	ID     string `json:"id"`
	Name   string `json:"name"`
}

func UnmarshalNewProject(data []byte) (NewProject, error) {
	var r NewProject
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *NewProject) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func InitProject(c *gin.Context) {
	resp, _ := c.GetRawData()
	temp := make(map[string]string)
	json.Unmarshal(resp, &temp)
	recordId := temp["record_id"]
	data := global.Cli.GetRecordInByte(P.AppToken, P.TableId, recordId)
	err := json.Unmarshal(data, &MyProject)
	if err != nil {
		panic(err)
	}
	//global.Cli.Send()
	if UserAccessToken != "" {
		CreateProject()
	}

}

func CreateProject() bool {
	var result bool = false
	if UserAccessToken == "" {
		err := errors.New("UserAccessToken为空，请再次鉴权！")
		panic(err)
		return result
	}
	pjt := MyProject.Data.Record.Fields
	var members []string
	for _, value := range pjt.ParticipatingMembers {
		members = append(members, value.ID)
	}
	manager := pjt.ProjectManager[0].ID

	chatId := global.Cli.CreateGroup("【"+pjt.ProjectProperties+"】"+pjt.ProjectName, "open_id", manager).ChatId
	global.Cli.AddMembers(chatId, "open_id", "1", members)
	spaceId := global.Cli.CreateKnowledgeSpace("【"+pjt.ProjectSource+"】"+pjt.ProjectName, pjt.ProjectProfile, "Bearer "+UserAccessToken).SpaceId
	robotId := global.Cli.GetRobotInfo().OpenId
	logrus.Println(robotId, "rob")
	global.Cli.AddMembersToKnowledgeSpace(spaceId, members, "open_id")

	result = true

	return result
}
