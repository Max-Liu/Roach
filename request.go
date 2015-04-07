package roach

import (
	"net"
	"net/http"
	"time"
)

type Request struct {
	*http.Request
	cookie string
}

func NewRequest(method string, url string, header map[string]string) (*http.Request, error) {
	//transport := http.Transport{
	//Dial: dialTimeout,
	//}

	newRequest, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, err
	}
	for k, v := range header {
		newRequest.Header.Add(k, v)
	}
	newRequest.Header.Add("Accept", `text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8`)
	newRequest.Header.Add("Accept-Encoding", `gzip,deflate,sdch"`)
	newRequest.Header.Add("Accept-Language", `en-US,en;q=0.8,ja;q=0.6,zh-CN;q=0.4,zh-TW;q=0.2`)
	newRequest.Header.Add("Connection", `keep-alive"`)
	newRequest.Header.Add("Cache-Control", "max-age=0")
	//cookieString := link.Request.cookie

	//req.Header.Add("Cookie", cookieString)
	newRequest.Header.Add("User-Agent", `mozila/5.0 (Macintosh; Intel Mac OS X 10_9_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36`)
	//newRequest.Client = &http.Client{Transport: &transport}
	return newRequest, nil
}
func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, time.Duration(requestTimeOut))
}
