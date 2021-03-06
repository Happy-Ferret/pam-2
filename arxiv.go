package main

import (
	"fmt"
	"net/http"
)

const (
	ArxivAPI = "http://export.arxiv.org/api"
)

func search(field, query string, maxResults int) error {
	url := fmt.Sprintf("%s/query?search_query=%s:\"%s\"&max_results=%d", ArxivAPI, field, query, maxResults)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	return nil
}

func SearchAll(query string, maxResults int) error {
	return nil
}

func SearchTitle(query string, maxResults int) error {
	return nil
}

func SearchAuthor(query string, maxResults int) error {
	return nil
}
