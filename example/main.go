package main

import (
	"roach"
	"time"

	"github.com/astaxie/beego/logs"
)

func main() {
	beeLog := logs.NewLogger(10000)
	beeLog.SetLogger("console", "")

	customHeader := make(map[string]string)
	customHeader["Accept"] = `text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8`
	customHeader["Accept-Encoding"] = `gzip,deflate,sdch`
	customHeader["Accept-Language"] = `en-US,en;q=0.8,ja;q=0.6,zh-CN;q=0.4,zh-TW;q=0.2`
	customHeader["Connection"] = `keep-alive`
	customHeader["Cache-Control"] = "max-age=0"

	config := &roach.ClanConfig{
		StartPoint:       "http://tieba.baidu.com",
		Host:             "baidu.com",
		ConcurrentNumber: 2,
		Rate:             500 * time.Millisecond,
		Log:              beeLog,
		Header:           customHeader,
		LinkConfig:       &roach.LinkConfig{},
	}

	clan := roach.NewClan(config)
	clan.Rush()
}
