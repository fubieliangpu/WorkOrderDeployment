package impl_test

import (
	"context"
	"crypto/md5"
	"fmt"
	"log"
	"testing"

	"github.com/fubieliangpu/WorkOrderDeployment/apps/user"
	"github.com/fubieliangpu/WorkOrderDeployment/ioc"
	"github.com/fubieliangpu/WorkOrderDeployment/test"
	"golang.org/x/crypto/bcrypt"
)

var (
	serviceImpl user.Service
	ctx         = context.Background()
)

func init() {
	test.DevelopmentSetup()

	serviceImpl = ioc.Controller.Get(user.AppName).(user.Service)
}

func TestCreateVisitorUser(t *testing.T) {
	req := user.NewCreateUserRuquest()
	req.Username = "visitor2"
	req.Password = "123456"
	req.Role = user.ROLE_VISITOR
	ins, err := serviceImpl.CreateUser(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ins)
}

func TestCreateAuthorUser(t *testing.T) {
	req := user.NewCreateUserRuquest()
	req.Username = "Admin1"
	req.Password = "123456"
	req.Role = user.ROLE_ADMIN
	ins, err := serviceImpl.CreateUser(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ins)
}

func TestQueryUser(t *testing.T) {
	req := user.NewQueryUserRequest()
	ins, err := serviceImpl.QueryUser(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ins)
}

func TestMd5(t *testing.T) {
	h := md5.New()
	h.Write([]byte("123456"))
	fmt.Printf("%x", h.Sum(nil))
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func TestUserCheckPassword(t *testing.T) {
	req := user.NewCreateUserRuquest()
	req.Username = "admin"
	req.Password = "123456"
	u := user.NewUser(req)
	u.HashPassword()
	t.Log(u.CheckPassword("123456"))
}

func TestPasswordHash(t *testing.T) {
	password := "123456"
	//hash, _ := HashPassword(password)
	hash := "$2a$10$QLIaQAX2iTx/RJe/sMRwZ.ZwTZR7HjZYrYHeIxJ.BIxoJpPxFB2Sa"
	fmt.Println("Password", password)
	fmt.Println("Hash:	", hash)
	match := CheckPasswordHash(password, hash)
	fmt.Println("Match:	", match)
}

// crnsm5914uficu1ve39g
func TestDeleteUser(t *testing.T) {
	req := user.NewDeleteUserRequest()
	req.AccessToken = "crntt9114ufmu00pggs0"
	req.Username = "Admin1"
	ins, err := serviceImpl.DeleteUser(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	t.Log(ins)
}

func TestChangeUser(t *testing.T) {
	req := user.NewChangeUserRequest()
	req.AccessToken = "crntt9114ufmu00pggs0"
	req.Username = "visitor2"
	req.Password = "43211"
	ins, err := serviceImpl.ChangeUser(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	t.Log(ins)
}
