package controller

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	global "xlab-feishu-robot/global"

	"github.com/gin-gonic/gin"
)

var (
	P         FeishuProjectFormPath
	T         TemplateDocs
	MyProject NewProject
)

// 立项问卷信息，从config读取
type FeishuProjectFormPath struct {
	AppToken string
	TableId  string
}

type TemplateDocs struct {
	SpaceId         string
	ParentNodeToken string
}
type Project struct {
	ProjectId int

	// project info
	ProjectName      string
	ProjectType      string // internal | external
	ProjectLeaderIds string // JSON: Array<int> (array of employeeIds)
	GroupId          string //不储存 groupId，因为一个 project 可能对应多个 group

	// doc related info:
	GanttDocUrl   string // 甘特图，其中的排期自动映射到飞书任务
	PrdDocUrl     string // PRD
	TechDocUrl    string // 技术文档
	FeishuRepoUrl string // 飞书知识空间首页

	// status
	ProjectStatus string // beforeStart, pending, revising, hang, finished, aborted
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

//先向用户发送鉴权链接，等待获取到UserAccessToken，然后再推送立项问卷链接。

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

	//if UserAccessToken != "" {
	//	CreateProject()
	//}

}

func setProjectType(projectProperties string) string {
	if strings.Contains(projectProperties, "内部") {
		return "internal"
	} else {
		return "external"
	}
}

func getProjectLeaderIds(p ParticipatingMember) string {
	var ProjectLeaderId interface{}
	//TODO:
	//在数据库中查找LeaderId

	//ProjectLeaderId=QueryEmployeeByFullname(p.Name) (*[]Employee, error)
	return "[" + "\"" + ProjectLeaderId.(string) + "\"" + "]"

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

	v := global.Cli.CreateGroup("【"+pjt.ProjectProperties+"】"+pjt.ProjectName, "open_id", manager)
	//将新建的群插入数据库，
	//InsertGroupRecords(v []Group, "open_id")
	global.Cli.AddMembers(v.ChatId, "open_id", "1", members)

	//创建知识空间
	s := global.Cli.CreateKnowledgeSpace("【"+pjt.ProjectSource+"】"+pjt.ProjectName, pjt.ProjectProfile, "Bearer "+UserAccessToken)

	var botIds []string
	robotId := global.Cli.GetRobotInfo().OpenId
	botIds = append(botIds, robotId)
	//设置群成员可见
	var chats []string
	chats = append(chats, v.ChatId)
	global.Cli.AddMembersToKnowledgeSpace(s.SpaceId, chats, "openchat")
	//将机器人设为管理员
	global.Cli.AddBotsToKnowledgeSpaceAsAdmin(s.SpaceId, botIds, UserAccessToken)

	//复制节点（生成原始文档）
	//需配置模板文档所在路径
	nodes := global.Cli.GetAllNodes(T.SpaceId, T.ParentNodeToken)
	for _, value := range nodes {
		global.Cli.CopyNode(T.SpaceId, value.NodeToken, s.SpaceId, "", "")
	}

	//插入项目信息至数据库
	/*
		projectInfo := Project{
			ProjectId:        0,
			ProjectName:      pjt.ProjectName,
			ProjectType:      setProjectType(pjt.ProjectProperties),
			ProjectLeaderIds: getProjectLeaderIds(pjt.ProjectManager[0]),
			//GroupId:          GetNewGroupId(),
			GanttDocUrl:   "",
			PrdDocUrl:     "",
			TechDocUrl:    "",
			FeishuRepoUrl: "",
			ProjectStatus: "",
		}
		var projectInfoList []Project
		projectInfoList = append(projectInfoList, projectInfo)
		InsertProjectRecords(projectInfoList)
	*/

	result = true
	//清除变量，为下一次立项准备
	UserAccessToken = ""
	p := reflect.ValueOf(MyProject).Elem()
	p.Set(reflect.Zero(p.Type()))

	return result
}
