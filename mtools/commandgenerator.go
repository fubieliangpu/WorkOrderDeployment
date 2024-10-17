package mtools

import (
	"io"
	"log"
	"os"
)

// 定义一个命令生成器，commandfile是命令写入的文件，可以被SSHTools使用
func CommandGenerator(commandfile string, commands ...string) {
	f, err := os.OpenFile(commandfile, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	for _, command := range commands {
		_, err := io.WriteString(f, command)
		if err != nil {
			log.Fatal(err)
		}
	}
}
