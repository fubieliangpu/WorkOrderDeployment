package user

import (
	"encoding/json"
	"fmt"

	"github.com/fubieliangpu/WorkOrderDeployment/common"
	"golang.org/x/crypto/bcrypt"
)

// 创建用户的参数
type CreateUserRequest struct {
	Username string `json:"username" validate:"required" gorm:"column:username"`
	Password string `json:"password" validate:"required" gorm:"column:password"`
	Role     Role   `json:"role" gorm:"column:role"`
	//将用户标签label以json格式存入数据库中，不需要在数据库中专门设计label id key value字段
	Label map[string]string `json:"label" gorm:"column:label;serializer:json"`
}

func NewCreateUserRuquest() *CreateUserRequest {
	return &CreateUserRequest{
		Role:  ROLE_VISITOR, //默认是普通用户
		Label: map[string]string{},
	}
}

func (req *CreateUserRequest) Validate() error {
	if req.Username == "" {
		return fmt.Errorf("please input username")
	}
	return nil
}

// 密码处理，hash
func (req *CreateUserRequest) HashPassword() error {
	cryptoPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	req.Password = string(cryptoPass)
	return nil
}

// 比较传来密码的hash值和已存密码的hash值
func (req *CreateUserRequest) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(req.Password), []byte(password))
}

// 用户对象创建
type User struct {
	*common.UserMeta
	*CreateUserRequest
}

func NewUser(req *CreateUserRequest) *User {
	return &User{
		UserMeta:          common.NewUserMeta(),
		CreateUserRequest: req,
	}
}

func (req *User) String() string {
	dj, _ := json.MarshalIndent(req, "", "	")
	return string(dj)
}

func (req *User) TableName() string {
	return "users"
}

// 系统用户列表集合
type UserSet struct {
	//总共有多少个
	Total int64 `json:"total"`
	//对象清单
	Items []*User `json:"items"`
}

func NewUserSet() *UserSet {
	return &UserSet{
		Items: []*User{},
	}
}
