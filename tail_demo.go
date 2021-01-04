package main

import (
	"fmt"
	"time"

	"github.com/hpcloud/tail"
)

//在日志文件中写入时要以回车键结尾才能读取出来
func main() {
	fileName := "./my.log"

	config := tail.Config{
		ReOpen:    true,                                 //重新打开
		Follow:    true,                                 //是否跟随
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, //从文件哪个地方读
		MustExist: false,                                //文件不存在不报错
		Poll:      true,
	}
	tails, err := tail.TailFile(fileName, config)
	if err != nil {
		fmt.Println("tail file failed,err:", err)
		return
	}
	var (
		line *tail.Line
		ok   bool
	)
	for {
		line, ok = <-tails.Lines
		if !ok {
			fmt.Printf("tail file close repon,fileName:%s\n", tails.Filename)
			time.Sleep(time.Second)
			continue
		}
		fmt.Println("Line:", line.Text)
	}
}
