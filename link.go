package roach

import (
	"bytes"
	"compress/gzip"
	"errors"
	"net"
	"net/http"
	"strings"
	"time"
)

var LinkStrChan chan string
var intChan chan int
var LinkChan chan *Link
var exitChan chan bool
var cookieChan chan string
var PureStack map[string]*Link

func init() {
	intChan = make(chan int)
	exitChan = make(chan bool)
	LinkChan = make(chan *Link, 10000)
	LinkStrChan = make(chan string, 10000)
	cookieChan = make(chan string, 1)
	PureStack = make(map[string]*Link)
}

type Link struct {
	Url          string
	Request      *Request
	config       *LinkConfig
	StatusCode   int
	duration     string
	error_count  int
	HasRequested bool
	HasGetUrl    bool
	Title        string
	log          Logger
}

type Request struct {
	*http.Client
	cookie string
}

type Roach struct {
}

func newLink(target string) *Link {
	newLink := &Link{
		Url:     target,
		config:  DefaultLinkConfigs,
		Request: NewRequest(),
	}
	return newLink
}

func (link *Link) setConfig(config *LinkConfig) {
	link.config = config
	if config.Log == nil {
		config.Log = DefaultLinkConfigs.Log
		config.Log.SetLogger("console", "")
	}
	link.log = config.Log
}

func NewRequest() *Request {
	transport := http.Transport{
		Dial: dialTimeout,
	}

	newRequest := &Request{
		cookie: ``,
	}
	newRequest.Client = &http.Client{
		Transport: &transport,
	}
	return newRequest
}

func (link *Link) MakeRequest() (*http.Response, error) {
	if link.HasRequested == true {
		return nil, errors.New("has Requested")
	}

	req, err := http.NewRequest("GET", string(link.Url), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", `text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8`)
	req.Header.Add("Accept-Encoding", `gzip,deflate,sdch"`)
	req.Header.Add("Accept-Language", `en-US,en;q=0.8,ja;q=0.6,zh-CN;q=0.4,zh-TW;q=0.2`)
	req.Header.Add("Connection", `keep-alive"`)
	req.Header.Add("Cache-Control", "max-age=0")
	cookieString := link.Request.cookie

	req.Header.Add("Cookie", cookieString)
	req.Header.Add("User-Agent", `mozila/5.0 (Macintosh; Intel Mac OS X 10_9_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36`)
	var breakCounter int
	var t0, t1 time.Time

	for {
		t0 = time.Now()
		link.log.Informational("Requesting %s", req.URL.String())
		resp, err := link.Request.Do(req)
		t1 = time.Now()
		if err != nil {
			breakCounter++
			if breakCounter == link.config.BadLinkRetryTimes {
				link.log.Warning("tried %s 3 times,abandoned")
				return nil, err
			}

			link.log.Warning("Request %s Response err: %s,retury in 5 Second", link.Url, err.Error())
			<-time.After(5 * time.Second)

		} else {
			link.StatusCode = resp.StatusCode
			if link.StatusCode != 200 {
				link.log.Warning("Request %s Response code:%d", link.Url, link.StatusCode)
			}

			link.duration = t1.Sub(t0).String()
			link.error_count = breakCounter
			link.HasRequested = true
			return resp, nil
		}
	}
}

func (l *Link) GetPageUrls() error {
	if l.HasGetUrl {
		return errors.New("Has Searched for the url")
	}
	res, err := l.MakeRequest()
	if err != nil {
		return err
	}

	var respBuffer bytes.Buffer
	respBuffer.ReadFrom(res.Body)
	gzipByte := respBuffer.Bytes()
	htmlByte := respBuffer.Bytes()

	gzipBuffer := bytes.NewBuffer(gzipByte)
	htmlBuffer := bytes.NewBuffer(htmlByte)

	gzipReader, err := gzip.NewReader(gzipBuffer)

	var contentStr string
	if err == nil {
		gzipBuffer.ReadFrom(gzipReader)
		contentStr = gzipBuffer.String()
		defer gzipReader.Close()
	} else {
		contentStr = htmlBuffer.String()
	}
	defer res.Body.Close()
	get_title(l, contentStr)

lookingForLink:
	for {
		startIndex := strings.Index(contentStr, `href="`)
		//check if it has looked entire page for href=
		if startIndex == -1 {
			break
		}

		newStr := contentStr[startIndex+6:]
		newStrEndIndex := strings.Index(newStr, `"`)

		if newStrEndIndex <= 6 {
			contentStr = newStr[2:]
			continue lookingForLink
		}
		linkStr := newStr[:newStrEndIndex]
		for _, v := range l.config.IgnoredFileExtention {
			if linkStr[len(linkStr)-4:len(linkStr)] == v {
				contentStr = newStr[newStrEndIndex:]
				continue lookingForLink
			}
		}

		if string(linkStr[0]) == "/" && string(linkStr[1]) == "/" {
			contentStr = newStr[newStrEndIndex+2:]
			continue lookingForLink
		}

		//check if linkSts is relative path.if so,change to absolute path.
		if string(linkStr[0]) == "/" {
			contentStr = newStr[newStrEndIndex:]

			LinkStrChan <- `http://` + l.config.Host + linkStr
			continue lookingForLink
		}

		if index := strings.Index(linkStr, l.config.Host); index == -1 {
			contentStr = newStr[newStrEndIndex:]
			continue lookingForLink
		}
		LinkStrChan <- linkStr
		contentStr = newStr[newStrEndIndex:]
	}
	l.HasGetUrl = true
	return nil
}

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, time.Duration(requestTimeOut))
}

func get_title(link *Link, contentStr string) {
	startIndex := strings.Index(contentStr, "<title>")
	endIndex := strings.Index(contentStr, "</title>")

	if startIndex != -1 && endIndex != -1 {
		link.Title = contentStr[startIndex+7 : endIndex]
	}
}
