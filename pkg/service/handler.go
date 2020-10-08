package service

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
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
	body, err := readRequest(r)
	if err != nil {
		httpError(w, http.StatusBadRequest)
		return
	}

	response := make([]payload, len(body))

	for i, url := range body {
		payload, err := visitURL(r.Context(), url)
		if err != nil {
			httpError(w, http.StatusBadRequest)
			return
		}
		response[i] = *payload
	}

	sendResponse(w, response)
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
	resp, err := http.Get(url)
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
