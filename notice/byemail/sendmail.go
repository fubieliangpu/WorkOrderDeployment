package byemail

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/smtp"
	"os"

	"github.com/jordan-wright/email"
)

// 把需要认证的信息放在结构体中
type MailAuthInfo struct {
	Identity      string
	SendAddress   string
	AuthorizeCode string
	ServerHost    string
	ServerPort    string
}

func NewMailAuthInfo() *MailAuthInfo {
	return &MailAuthInfo{}
}

func (m *MailAuthInfo) String() string {
	dj, _ := json.MarshalIndent(m, "", "	")
	return string(dj)
}

// 自定义的邮件发送方法
func (m *MailAuthInfo) MySendmail(e *email.Email) error {
	if err := e.Send(
		fmt.Sprintf("%v:%v", m.ServerHost, m.ServerPort),
		smtp.PlainAuth(m.Identity, m.SendAddress, m.AuthorizeCode, m.ServerHost),
	); err != nil {
		return err
	}
	return nil
}

// 自定义NewEmail方法，基于原NewEmail,主要是将模版生成后的HTML的内容赋给email.Email.HTML
func MyNewEmail(filename string) *email.Email {
	return &email.Email{
		HTML: *func(fname string) *[]byte {
			f, _ := os.Open(fname)
			defer f.Close()
			content, err := io.ReadAll(f)
			if err != nil {
				log.Fatal(err)
			}
			return &content
		}(filename),
	}
}
