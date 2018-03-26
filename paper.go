package main

import "net/url"

type Paper struct {
	Title    string   `json:"title"`
	Authors  []string `json:"authors"`
	Abstract string   `json:"abstract"`
	Note     string   `json:"note"`
	Favorite bool     `json:"favorite"`
	Read     bool     `json:"read"`
	Master   bool     `json:"master"`
	Tags     []string `json:"tags"`
	URL      url.URL  `json:"url"`
}
