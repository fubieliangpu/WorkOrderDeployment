package token_test

import (
	"testing"
	"time"

	"github.com/fubieliangpu/WorkOrderDeployment/apps/token"
	"github.com/fubieliangpu/WorkOrderDeployment/apps/user"
)

func TestTokenString(t *testing.T) {
	tk := token.Token{
		UserId: 1,
		Role:   user.ROLE_VISITOR,
	}

	t.Log(tk.String())
}

func TestNewToken(t *testing.T) {
	req := user.NewCreateUserRuquest()
	req.Username = "fublp"
	req.Password = "324134"
	req.Role = user.ROLE_ADMIN
	nus := user.NewUser(req)
	tk := token.NewToken(nus)
	t.Log(tk)
}

func TestTokenExpired(t *testing.T) {
	nowtime := time.Now().Unix()
	tk := token.Token{
		UserId:               1,
		Role:                 user.ROLE_ADMIN,
		AccessTokenExpiredAt: 5,
		CreatedAt:            nowtime,
	}
	t.Log(tk.AccessTokenIsExpired())
}
