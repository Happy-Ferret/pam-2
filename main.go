package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

// ~/.pam stores all metadata and PDF documents for the papers.
var (
	path      string
	indexTmpl string
	paperTmpl string
	pam       *Pam
)

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
	// Load index.tmpl
	tmpl, err := ioutil.ReadFile("resources/index.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	indexTmpl = string(tmpl)
	// Load paper.tmpl
	tmpl, err = ioutil.ReadFile("resources/paper.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	paperTmpl = string(tmpl)
	// Create a new Pam
	pam, err = NewPam()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	fs := http.FileServer(http.Dir("./resources/"))
	http.Handle("/resources/", http.StripPrefix("/resources/", fs))
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/p/", paperHandler)
	go openWeb("http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, err := template.New("index").Funcs(template.FuncMap{
			"join": strings.Join,
		}).Parse(indexTmpl)
		if err != nil {
			log.Fatal(err)
		}
		if err = pam.Reload(); err != nil {
			log.Fatal(err)
		}
		pam.Reload()
		t.Execute(w, pam)
	} else if r.Method == "POST" {
		r.ParseForm()
		fmt.Println(r.Form["search"])
	}
}

func paperHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("paper").Funcs(template.FuncMap{
		"join": strings.Join,
		"checked": func(checked bool) string {
			if checked {
				return "checked"
			}
			return ""
		},
	}).Parse(paperTmpl)
	if err != nil {
		log.Fatal(err)
	}
	title := r.URL.Path[len("/p/"):]
	p := pam.Papers.SearchByTitle(title)
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
