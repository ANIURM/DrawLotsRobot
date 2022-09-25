package model

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DepartmentType string

const (
	Software  DepartmentType = "software"
	Hardward  DepartmentType = "hardware"
	Operation DepartmentType = "operation"
	Design    DepartmentType = "design"
	Media     DepartmentType = "media"
)

type RoleType string

const (
	Member   RoleType = "member"
	Intern   RoleType = "intern"
	Resigned RoleType = "resigned"
)

type Employee struct {
	EmployeeId primitive.ObjectID `bson:"_id,omitempty"`

	FeishuOpenId string `bson:"feishu_openid,omitempty"`
	FeishuUserId string `bson:"feishu_userid,omitempty"`

	Fullname    string             `bson:"fullname,omitempty"`
	School      string             `bson:"school,omitempty"`
	StudentId   string             `bson:"student_id"`
	CitizenId   string             `bson:"citizen_id,omitempty"`
	Phone       string             `bson:"phone,omitempty"`
	ClaimGender string             `bson:"claim_gender"`
	Birthday    primitive.DateTime `bson:"birthday,omitempty"`

	Department        DepartmentType     `bson:"department,omitempty"`
	Role              RoleType           `bson:"role,omitempty"`
	PrivilegeGroupIds string             `bson:"privilege_group_ids,omitempty"`
	InternStartDate   primitive.DateTime `bson:"intern_start_date,omitempty"`
	MemberStartDate   primitive.DateTime `bson:"member_start_date"`
	ResignDate        primitive.DateTime `bson:"resign_date"`
}

func InsertEmployeeRecords(v []Employee) {
	var l []interface{}

	for _, value := range v {
		l = append(l, value)
	}

	employee.InsertMany(context.TODO(), l)
}

func QueryEmployeeByFullname(name string) (*[]Employee, error) {
	var employees []Employee

	filter := bson.D{{Key: "fullname", Value: name}}

	cur, err := employee.Find(context.TODO(), filter)
	if err != nil {
		logrus.New().WithField("ProjectStatus", name).Error(err)
		return nil, err
	}

	var result Employee
	for cur.Next(context.TODO()) {
		cur.Decode(&result)
		employees = append(employees, result)
	}
	return &employees, nil
}

func QueryEmployeeByStudentId(sid string) (*[]Employee, error) {
	var employees []Employee

	filter := bson.D{{Key: "student_id", Value: sid}}

	cur, err := employee.Find(context.TODO(), filter)
	if err != nil {
		logrus.New().WithField("StudentId", sid).Error(err)
		return nil, err
	}

	var result Employee
	for cur.Next(context.TODO()) {
		cur.Decode(&result)
		employees = append(employees, result)
	}
	return &employees, nil
}

func QueryEmployeeByCitizenId(cid string) (*[]Employee, error) {
	var employees []Employee

	filter := bson.D{{Key: "citizen_id", Value: cid}}

	cur, err := employee.Find(context.TODO(), filter)
	if err != nil {
		logrus.New().WithField("CitizenId", cid).Error(err)
		return nil, err
	}

	var result Employee
	for cur.Next(context.TODO()) {
		cur.Decode(&result)
		employees = append(employees, result)
	}
	return &employees, nil
}

func QueryEmployeeByDepartment(department DepartmentType) (*[]Employee, error) {
	var employees []Employee

	filter := bson.D{{Key: "department", Value: department}}

	cur, err := employee.Find(context.TODO(), filter)
	if err != nil {
		logrus.New().WithField("Department", department).Error(err)
		return nil, err
	}

	var result Employee
	for cur.Next(context.TODO()) {
		cur.Decode(&result)
		employees = append(employees, result)
	}
	return &employees, nil
}

func QueryEmployeeByRole(role RoleType) (*[]Employee, error) {
	var employees []Employee

	filter := bson.D{{Key: "role", Value: role}}

	cur, err := employee.Find(context.TODO(), filter)
	if err != nil {
		logrus.New().WithField("Role", role).Error(err)
		return nil, err
	}

	var result Employee
	for cur.Next(context.TODO()) {
		cur.Decode(&result)
		employees = append(employees, result)
	}
	return &employees, nil
}
