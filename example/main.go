package main

import (
	"roach"
	"time"

	"github.com/astaxie/beego/logs"
)

func main() {
	beeLog := logs.NewLogger(10000)
	beeLog.SetLogger("console", "")

	config := &roach.ClanConfig{
		StartPoint:       "http://tieba.baidu.com",
		Host:             "baidu.com",
		ConcurrentNumber: 2,
		Rate:             1 * time.Second,
		Log:              beeLog,
	}

	clan := roach.NewClan(config)
	clan.Rush()
}
