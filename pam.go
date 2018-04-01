package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"path/filepath"
	"time"
)

var pamQuotes = []string{
	"Absolutely I do",
	"There's a lot of beauty in ordinary things. Isn't that kind of the point?",
	"I think you're a witch.",
	"He talked about himself in the third person?",
	"Why do you need to wear the holster at all?",
	"No laughing. No comments. Just positive energy and we'll have a pure fun day. Okay?",
	"So everyone here knows pirate code?",
	"Dwight! What are you doing? We've only been in here for like two seconds!",
	"You can flirt with someone to get what you want and also be attracted to them. How do you think we got together?",
	"Sometimes the heart doesn't know what it wants until it finds what it wants.",
	"Wanna count her fingers and toes again?",
	"I cannot wait for that joke to be over.",
	"Oscar and the warehouse guy! Go Oscar! Go gay warehouse guy!",
	"Don't do the twirl.",
}

type Pam struct {
	PamQuote string
	Papers   Papers
}

func NewPam() (*Pam, error) {
	papers, err := importPapers(path)
	if err != nil {
		return nil, err
	}
	papers.SortByTitle()
	return &Pam{
		PamQuote: randPamQuote(),
		Papers:   papers,
	}, nil
}

func (p *Pam) Reload() error {
	p.PamQuote = randPamQuote()
	papers, err := importPapers(path)
	if err != nil {
		return err
	}
	papers.SortByTitle()
	p.Papers = papers
	return nil
}

// randPamQuote returns a random Pam Quote.
func randPamQuote() string {
	rand.Seed(time.Now().UnixNano())
	return pamQuotes[rand.Intn(len(pamQuotes))]
}

// listPapers returns a slice of file names (only PDF) in the path.
func importPapers(dir string) (Papers, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	papers := []*Paper{}
	for _, f := range files {
		filename := f.Name()
		if filepath.Ext(filename) == ".json" {
			p, err := FromJSON(filepath.Join(dir, filename))
			if err != nil {
				return nil, err
			}
			papers = append(papers, p)
		}
	}
	return papers, nil
}
