package byemail

import (
	"encoding/json"
	"html/template"
	"os"
)

type Page struct {
	Title        string
	ContentTitle string
	Message      string
}

func NewPage() *Page {
	return &Page{}
}

func (p *Page) String() string {
	dj, _ := json.MarshalIndent(p, "", "	")
	return string(dj)
}

// 通过HTML模版文件生成邮件回复正文的HTML文件
func (p *Page) GenHTMLFromModel(modelfile, htmlfile string) error {
	tmpl, err := template.ParseFiles(modelfile)
	if err != nil {
		return err
	}
	f, _ := os.OpenFile(htmlfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	defer f.Close()

	err = tmpl.Execute(f, *p)
	if err != nil {
		return err
	}
	return nil
}
