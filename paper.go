package main

import (
	"encoding/json"
	"net/url"
	"os"
)

type Paper struct {
	Title    string   `json:"title"`
	Authors  []string `json:"authors"`
	Abstract string   `json:"abstract"`
	Note     string   `json:"note"`
	Favorite bool     `json:"favorite"`
	Read     bool     `json:"read"`
	Master   bool     `json:"master"`
	Tags     []string `json:"tags"`
	FilePath string   `json:"filepath"`
	URL      url.URL  `json:"url"`
}

// FromJSON creates a new Paper from the JSON file, given by the argument file name.
func FromJSON(filename string) (*Paper, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	paper := &Paper{}
	decoder := json.NewDecoder(f)
	if err = decoder.Decode(&paper); err != nil {
		return nil, err
	}
	return paper, nil
}
