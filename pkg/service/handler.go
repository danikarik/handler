package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	_defaultTimeout     = 1 * time.Second
	_defaultWorkerCount = 4
)

type urls []string

func (u urls) IsValid() bool {
	return u != nil && len(u) <= 20
}

type payload struct {
	URL      string `json:"url"`
	Response string `json:"response"`
}

func urlHandler(w http.ResponseWriter, r *http.Request) {
	urls, err := readRequest(r)
	if err != nil {
		httpError(w, http.StatusBadRequest)
		return
	}

	var (
		pool = make(chan struct{}, _defaultWorkerCount)
		done = make(chan payload, len(urls))
	)

	var (
		exit = make(chan struct{})
		errC = make(chan error)
	)

	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			pool <- struct{}{}

			// Visit page.
			payload, err := visitURL(r.Context(), url)
			if err != nil {
				errC <- err
			} else {
				done <- *payload
			}

			<-pool
		}(url)
	}

	go func() {
		wg.Wait()
		close(done)
		close(exit)
	}()

	select {
	case <-exit:
		response := make([]payload, 0)
		for resp := range done {
			response = append(response, resp)
		}
		sendResponse(w, response)
	case err := <-errC:
		log.Println(err)
		httpError(w, http.StatusBadRequest)
	case <-r.Context().Done():
		log.Println("context done")
		httpError(w, http.StatusBadRequest)
	}
}

func readRequest(r *http.Request) (urls, error) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var body urls
	if err := json.Unmarshal(data, &body); err != nil {
		return nil, err
	}

	if !body.IsValid() {
		return nil, errors.New("request body has invalid length")
	}

	return body, nil
}

func sendResponse(w http.ResponseWriter, r []payload) error {
	data, err := json.Marshal(r)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	return err
}

func visitURL(ctx context.Context, url string) (*payload, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	fmt.Println("visiting " + url)

	client := &http.Client{Timeout: _defaultTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return &payload{URL: url, Response: string(data)}, nil
}
