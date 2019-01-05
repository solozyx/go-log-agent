package main

import (
	"fmt"
	"time"
	"github.com/hpcloud/tail"
)

func main() {
	filename := "tail_demo.go"
	tails, err := tail.TailFile(filename, tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	})
	if err != nil {
		fmt.Println("tail open file err:", err)
		return
	}
	var msg *tail.Line
	var ok bool
	for true {
		msg, ok = <- tails.Lines
		if !ok {
			fmt.Printf("tail file close reopen, filename:%s\n", tails.Filename)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		fmt.Println("msg:", msg)
	}
}