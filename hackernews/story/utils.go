package story

import (
	"encoding/json"
	"io"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/theonejonahgold/pwa/hackernews"
	"github.com/theonejonahgold/pwa/hackernews/comment"
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

// GetTopStories gets 50 stories from the hackernews top stories api.
func GetTopStories() ([]hackernews.HackerNewsObject, error) {
	res, err := client.Get("https://hacker-news.firebaseio.com/v0/topstories.json")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var storyIDs StoryIDArray
	if err = json.Unmarshal(b, &storyIDs); err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	c := make(chan hackernews.HackerNewsObject, len(storyIDs))
	for _, value := range storyIDs {
		wg.Add(1)
		go storyWorker(strconv.Itoa(value), c, &wg)
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	stories := make([]hackernews.HackerNewsObject, 0, len(storyIDs))
	for s := range c {
		if s == nil || s.GetType() != "story" {
			continue
		}
		stories = append(stories, s)
	}
	return stories, nil
}

func storyWorker(id string, c chan<- hackernews.HackerNewsObject, wg *sync.WaitGroup) {
	defer wg.Done()

	b, err := client.Get("https://hacker-news.firebaseio.com/v0/item/" + id + ".json")
	if err != nil {
		return
	}

	s, err := Parse(b)
	if err != nil {
		return
	}

	c <- s
}

func commentWorker(jc <-chan string, cc chan<- hackernews.HackerNewsObject, wg *sync.WaitGroup) {
	defer wg.Done()

	for id := range jc {
		res, err := client.Get("https://hacker-news.firebaseio.com/v0/item/" + id + ".json")
		if err != nil {
			continue
		}

		cm, err := comment.Parse(res)
		res.Body.Close()
		if err != nil {
			continue
		}
		cc <- cm
	}
}

func Parse(res *http.Response) (hackernews.HackerNewsObject, error) {
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	s := New()
	if err := json.Unmarshal(b, &s); err != nil {
		return nil, err
	}
	return s, nil
}
