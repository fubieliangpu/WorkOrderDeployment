package user_test

import (
	"fmt"
	"testing"

	"github.com/fubieliangpu/WorkOrderDeployment/apps/user"
)

func TestUserValidate(t *testing.T) {
	cu := user.NewCreateUserRuquest()
	cu.Username = "fublp"
	err := cu.Validate()
	t.Log(err)
}

func TestHashPassword(t *testing.T) {
	req := user.NewCreateUserRuquest()
	req.Password = "123456"
	req.HashPassword()
	fmt.Println(req)
}

func TestCheckPassword(t *testing.T) {
	req := user.NewCreateUserRuquest()
	req.Password = "11111"
	req.HashPassword()
	err := req.CheckPassword("1111")
	t.Log(err)
}

func TestStringUser(t *testing.T) {
	req := user.NewCreateUserRuquest()
	ins := user.NewUser(req)
	t.Log(ins)
}
