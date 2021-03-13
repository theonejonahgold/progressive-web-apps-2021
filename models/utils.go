package models

import (
	"encoding/json"
	"io"
	"net"
	"net/http"
	"sync"
	"time"
)

var (
	// Solution from: https://github.com/golang/go/issues/13998
	client = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			DisableKeepAlives:     true,
			IdleConnTimeout:       10 * time.Second,
		},
	}
)

func fetchComment(jc <-chan string, cc chan<- *Comment, wg *sync.WaitGroup) {
	defer wg.Done()
	for id := range jc {
		res, err := client.Get("https://hacker-news.firebaseio.com/v0/item/" + id + ".json")
		if err != nil {
			continue
		}

		bytes, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			continue
		}

		cm := NewComment()
		err = json.Unmarshal(bytes, cm)
		if err != nil {
			continue
		}
		cc <- cm
	}
}

func CloseWhenWGIsDone(c chan *Comment, wg *sync.WaitGroup) {
	wg.Wait()
	close(c)
}
