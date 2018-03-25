package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
)

// ~/.pam stores all metadata and PDF documents for the papers.
var path string

func init() {
	// Check if ~/.pam exists; create a new one if not.
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	path = filepath.Join(usr.HomeDir, ".pam")
	if _, err = os.Stat(path); os.IsNotExist(err) {
		if err = os.Mkdir(path, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	http.HandleFunc("/", mainHandler)
	go openWeb("http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	papers := listPapers(path)
	for _, p := range papers {
		fmt.Fprintf(w, "%s\n", p)
	}
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
