package tNet

import (
	"fmt"
)

var GlobalSysEventChan = make(chan SysEventValInterface)

type SysEventValInterface interface {
	GetFunc() func([]any)
	GetArgs() []any
}

// 异步消息处理
func init() {
	go func() {
		for {
			select {
			case event := <-GlobalSysEventChan:
				fun := event.GetFunc()
				args := event.GetArgs()
				go func(){
					defer func() {
						if r := recover(); r != nil {
							fmt.Println("PANIC_ERROR:", r)
						}
					}()
					fun(args)
				}()
			}
		}
	}()
}
