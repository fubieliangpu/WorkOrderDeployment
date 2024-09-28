package mtools

import (
	"io"
	"log"
	"os"
	"regexp"

	"github.com/fubieliangpu/WorkOrderDeployment/exception"
	"github.com/rs/zerolog"
)

// 定义一个冲突检测正则匹配工具
func Regexper(recordfile string, regset ...string) (bool, error) {
	//自定义匹配结果日志，如果匹配到则记录日志
	zerolog.TimeFieldFormat = "2006/01/02 15:04:05 -0700"
	flog, err := os.OpenFile("regexp.log", os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	defer flog.Close()

	//自定义logger
	logger := zerolog.New(flog).With().Timestamp().Logger()

	//打开登录设备进行配置或者信息查询的记录文件，用于正则表达式匹配
	f, err := os.Open(recordfile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	recordcontent, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	//记录正则匹配到的次数，如果所有正则没有匹配，则数量一定小于表达式的数量，则返回失败结果
	successcount := 0
	for _, singlereg := range regset {
		re := regexp.MustCompile(singlereg)
		for _, v := range re.FindAll(recordcontent, -1) {
			if len(v) > 0 {
				//记录日志文件，自定义日志
				logger.Info().Msgf("Match content '%v' by expression '%v'", string(v), singlereg)
				successcount++
			}
		}
	}
	if successcount == len(regset) {

		return true, nil
	}
	//nil后期替换为自定义错误类型
	return false, exception.ErrRegularMatchFailed("Regular expression matching failed, there may be conflicts!")
}
