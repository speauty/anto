package cron

import (
	"anto/lib/log"
	"anto/lib/util"
	"context"
	"fmt"
	"runtime"
)

func FuncSrtCronMsgRedirect(ctx context.Context, cronName string, ptrLog *log.Log, chanMsg, chanMsgRedirect chan string) {
	go func(ctx context.Context, localChanMsg, localChanMsgRedirect chan string) {
		coroutineName := "消息协程"
		chanName := "chanMsg"
		for true {
			select {
			case <-ctx.Done():
				ptrLog.WarnF("%s关闭(ctx.done), %s被迫退出", cronName, coroutineName)
				runtime.Goexit()
			case currentMsg, isOpen := <-localChanMsg:
				if isOpen == false && currentMsg == "" {
					ptrLog.WarnF("%s-%s通道关闭, %s被迫退出", cronName, chanName, coroutineName)
					runtime.Goexit()
				}
				if localChanMsgRedirect != nil {
					localChanMsgRedirect <- fmt.Sprintf("时间: %s, 来源: %s, 信息: %s", util.GetDateTime(), cronName, currentMsg)
				}
			}
		}
	}(ctx, chanMsg, chanMsgRedirect)
}
