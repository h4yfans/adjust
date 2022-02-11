package tool

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"

	client "github.com/h4yfans/adjust/http"
)

type Tool struct {
	parallelChan chan struct{}
	wg           sync.WaitGroup
	httpClient   client.HTTPClient
}

func NewTool(parallel uint, client client.HTTPClient) *Tool {
	return &Tool{
		parallelChan: make(chan struct{}, parallel),
		wg:           sync.WaitGroup{},
		httpClient:   client,
	}
}

// Run makes all the requests, prints the result. It blocks until
// all work is done.
func (t *Tool) Run(urls ...string) error {
	var validURLs []string

	for _, u := range urls {
		if validURL, err := t.parseRawURL(u); err != nil {
			return fmt.Errorf("could not parse %s, error: %v", u, err)
		} else {
			validURLs = append(validURLs, validURL.String())
		}
	}

	return t.run(validURLs...)
}

func (t *Tool) parseRawURL(rawURL string) (*url.URL, error) {
	u, err := url.Parse(rawURL)
	if err != nil || u.Scheme == "" {
		u, err = url.ParseRequestURI("http://" + rawURL)
		if err != nil {
			return nil, err
		}
	}

	if !strings.Contains(u.Host, ".") {
		return nil, errors.New("url do not contains host")
	}

	return u, nil
}

func (t *Tool) run(urls ...string) error {

	t.wg.Add(len(urls))

	for i := range urls {
		go func(i int) {
			t.parallelChan <- struct{}{}
			defer func() {
				t.wg.Done()
				<-t.parallelChan
			}()

			md5Hash, err := createMD5Hash(t.httpClient, urls[i])
			if err != nil {
				fmt.Printf("%s %s\n", urls[i], err)
			} else {
				fmt.Printf("%s %s\n", urls[i], hex.EncodeToString(md5Hash[:]))
			}
		}(i)
	}

	t.wg.Wait()
	return nil
}

func createMD5Hash(httpClient client.HTTPClient, url string) ([md5.Size]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return [md5.Size]byte{}, fmt.Errorf("failed to create new http request: %s, error: %v", url, err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return [md5.Size]byte{}, fmt.Errorf("failed to request: %s, error: %v", url, err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return [md5.Size]byte{}, fmt.Errorf("failed to read request's body: %s, error: %v", url, err)
	}

	return md5.Sum(body), nil
}
