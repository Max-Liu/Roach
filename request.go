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
	return newRequest, nil
}
func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, time.Duration(requestTimeOut))
}
