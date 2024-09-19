package user

//创建用户的参数
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

type User struct {
}
