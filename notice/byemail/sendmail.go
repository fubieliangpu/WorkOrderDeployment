package byemail

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/smtp"
	"os"

	"github.com/fubieliangpu/WorkOrderDeployment/exception"
	"github.com/jordan-wright/email"
	"gopkg.in/yaml.v3"
)

// 把需要认证的信息放在结构体中
type MailAuthInfo struct {
	Identity      string `toml:"identity" yaml:"identity" json:"identity"`
	SendAddress   string `toml:"sendaddr" yaml:"sendaddr" json:"sendaddr"`
	AuthorizeCode string `toml:"authorizecode" yaml:"authorizecode" json:"authorizecode"`
	ServerHost    string `toml:"serverhost" yaml:"serverhost" json:"serverhost"`
	ServerPort    string `toml:"serverport" yaml:"serverport" json:"serverport"`
	//补充收件信息
	RecvAddress string `toml:"revcaddr" yaml:"revcaddr" json:"revcaddr"`
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

// 从指定yaml文件中加载邮件认证信息
func (m *MailAuthInfo) LoadEmailAuthInfoFromYaml(filename string) (*MailAuthInfo, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, exception.ErrOpenFileFailed(err.Error())
	}

	err = yaml.Unmarshal(content, m)
	if err != nil {
		return nil, exception.ErrParseFileFailed(err.Error())
	}

	return m, nil
}
