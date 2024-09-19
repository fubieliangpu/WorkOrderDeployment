package user

//定义角色，管理员或普通用户

type Role int

const (
	ROLE_ADMIN Role = iota + 1
	ROLE_VISITOR
)
