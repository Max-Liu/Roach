package roach

import (
	"time"

	"github.com/astaxie/beego/logs"
)

var ignoredFileExtention = []string{".css", ".js", ".png", ".jpg", ".ico"}
var host = "baidu.com"
var target = "http://tieba.baidu.com"
var linkStack *[]Link
var badLinkRetryTimes = 2
var requestTimeOut = 500 * time.Millisecond
var logger = logs.NewLogger(1000)

type LinkConfig struct {
	IgnoredFileExtention []string
	Target               string
	Host                 string
	BadLinkRetryTimes    int
	RequestTimeOut       time.Duration
	Log                  Logger
	header               map[string]string
}

type ClanConfig struct {
	ConcurrentNumber int
	Rate             time.Duration
	LinkConfig       *LinkConfig
	StartPoint       string
	Host             string
	Log              Logger
	Header           map[string]string
}

var DefaultClanConfigs = &ClanConfig{
	ConcurrentNumber: 5,
	Rate:             time.Millisecond,
	LinkConfig:       DefaultLinkConfigs,
}

var DefaultLinkConfigs = &LinkConfig{
	IgnoredFileExtention: ignoredFileExtention,
	Host:                 host,
	Target:               target,
	BadLinkRetryTimes:    badLinkRetryTimes,
	RequestTimeOut:       requestTimeOut,
	Log:                  logger,
}
