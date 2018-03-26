package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
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
	"It's not about who you've been with. It's about who you end up with. Sometimes the heart doesn't know what it wants until it finds what it wants.",
	"Wanna count her fingers and toes again?",
	"I cannot wait for that joke to be over.",
	"Oscar and the warehouse guy! Go Oscar! Go gay warehouse guy!",
	"Don't do the twirl.",
}

// ~/.pam stores all metadata and PDF documents for the papers.
var path string

func init() {
	// Check if ~/.pam exists
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	path = filepath.Join(usr.HomeDir, ".pam")
	if _, err = os.Stat(path); os.IsNotExist(err) {
		log.Printf(".pam not found, creating %s...", path)
		if err = os.Mkdir(path, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}
	// Create library folder for document files
	libPath := filepath.Join(path, "library")
	if _, err = os.Stat(libPath); os.IsNotExist(err) {
		log.Printf("library not found in ~/.pam, creating %s...", libPath)
		if err = os.Mkdir(libPath, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	fs := http.FileServer(http.Dir("./resources/"))
	http.Handle("/resources/", http.StripPrefix("/resources/", fs))
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/paper/", paperHandler)
	go openWeb("http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := ioutil.ReadFile("resources/index.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	t, err := template.New("index").Parse(string(tmpl))
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(w, NewPam())
}

func paperHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := ioutil.ReadFile("resources/paper.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	t, err := template.New("paper").Funcs(template.FuncMap{
		"join": strings.Join,
		"checked": func(checked bool) string {
			if checked {
				return "checked"
			}
			return ""
		},
	}).Parse(string(tmpl))
	if err != nil {
		log.Fatal(err)
	}

	// TODO: get paper name from response

	// Placeholder for visuals
	p := &Paper{
		Title: "Convolution by Evolution: Differentiable Pattern Producing Networks",
		Authors: []string{"Chrisantha Fernando", "Dylan Banarse", "Malcolm Reynolds",
			"Frederic Besse", "David Pfau", "Max Jaderberg", "Marc Lanctot", "Daan Wierstra"},
		Abstract: "In this work we introduce a differentiable version of the Compositional Pattern Producing Network, called the DPPN. Unlike a standard CPPN, the topology of a DPPN is evolved but the weights are learned. A Lamarckian algorithm, that combines evolution and learning, produces DPPNs to reconstruct an image. Our main result is that DPPNs can be evolved/trained to compress the weights of a denoising autoencoder from 157684 to roughly 200 parameters, while achieving a reconstruction accuracy comparable to a fully connected network with more than two orders of magnitude more parameters. The regularization ability of the DPPN allows it to rediscover (approximate) convolutional network architectures embedded within a fully connected architecture. Such convolutional architectures are the current state of the art for many computer vision applications, so it is satisfying that DPPNs are capable of discovering this structure rather than having to build it in by design. DPPNs exhibit better generalization when tested on the Omniglot dataset after being trained on MNIST, than directly encoded fully connected autoencoders. DPPNs are therefore a new framework for integrating learning and evolution. ",
		Note:     "This is my note for this paper!",
		Favorite: true,
		Read:     false,
		Master:   false,
	}

	t.Execute(w, p)
}

// openWeb executes a command that opens the default browser with the argument URL.
func openWeb(url string) {
	var cmd string
	var args []string
	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	if err := exec.Command(cmd, args...).Start(); err != nil {
		log.Fatal(err)
	}
}

type Pam struct {
	PamQuote string
	Titles   []string
}

func NewPam() *Pam {
	return &Pam{
		PamQuote: pamQuotes[rand.Intn(len(pamQuotes))],
		Titles:   listPapers(path),
	}
}

// listPapers returns a slice of file names (only PDF) in the path.
func listPapers(dir string) []string {
	// TODO: check/create metadata for performance
	// TODO: list the papers in order of last opened; alphabetical otherwise.

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	list := []string{}
	for _, f := range files {
		filename := f.Name()
		if filepath.Ext(filename) == ".pdf" {
			list = append(list, filename)
		}
	}
	return list
}
