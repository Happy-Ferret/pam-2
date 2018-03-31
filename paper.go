package main

import (
	"encoding/json"
	"net/url"
	"os"
	"sort"
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

type byTitle []*Paper

func (t byTitle) Len() int           { return len(t) }
func (t byTitle) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t byTitle) Less(i, j int) bool { return t[i].Title < t[j].Title }

type byAuthor []*Paper

func (a byAuthor) Len() int           { return len(a) }
func (a byAuthor) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byAuthor) Less(i, j int) bool { return a[i].Authors[0] < a[j].Authors[0] }

type Papers []*Paper

// SortByTitle sorts the list of paper by titles in the alphabetical order.
func (p Papers) SortByTitle() {
	sort.Sort(byTitle(p))
}

// SortByAuthor sorts the list of paper by the first author in the alphabetical order.
func (p Papers) SortByAuthor() {
	sort.Sort(byAuthor(p))
}
