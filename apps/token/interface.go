package token

import (
	"context"
)

//APPname

const (
	AppName = "token"
)

type Service interface {
	//令牌颁发
	IssueToken(context.Context, *IssueTokenRequest) (*Token, error)
	//令牌撤销
	RevolkToken(context.Context, *RevolkTokenRequest) (*Token, error)
	//令牌校验
	ValidateToken(context.Context, *ValidateTokenRequest) (*Token, error)
}

type IssueTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	IsMember bool   `json:"is_member"`
}

func NewIssueTokenRequest(username, password string) *IssueTokenRequest {
	return &IssueTokenRequest{
		Username: username,
		Password: password,
		IsMember: false,
	}
}

type RevolkTokenRequest struct {
	AccessToken  string
	RefreshToken string
}

func NewRevolkTokenRequest(at, rt string) *RevolkTokenRequest {
	return &RevolkTokenRequest{
		AccessToken:  at,
		RefreshToken: rt,
	}
}

type ValidateTokenRequest struct {
	AccessToken string
}

func NewValidateTokenRequest(at string) *ValidateTokenRequest {
	return &ValidateTokenRequest{
		AccessToken: at,
	}
}
