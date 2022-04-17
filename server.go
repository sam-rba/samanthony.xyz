package main

import (
	"fmt"
	tmpl "html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	fp "path/filepath"
	"strings"
)

const (
	host = ""
	port = "6969"
	root = "htdocs/"
)

var tmpls = make(map[string]*tmpl.Template)

func init() {
	err := fp.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if fp.Clean(path) == fp.Clean(root) ||
			fp.Ext(path) != ".html" ||
			path == fp.Join(root, "base.html") {
			return nil
		}
		label := path[len(fp.Clean(root)):]
		tmpls[label] = tmpl.Must(tmpl.ParseFiles(fp.Join(root, "base.html"), path))
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

type Page struct {
	Nav Nav
}

type Nav struct {
	ThisSection string
	Links       []NavLink
}

type NavLink struct {
	Href  string
	Label string
}

var nav = Nav{
	Links: []NavLink{
		{"/", "samanthony.xyz"},
		{"/software/", "software"},
	},
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	reqPath := r.URL.Path

	// If request directory, serve index.html.
	// ie. /software -> /software/index.html
	if info, err := os.Stat(fp.Join(root, reqPath)); err == nil {
		if info.IsDir() {
			reqPath = path.Join(reqPath, "index.html")
		}
	} else if os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	} else {
		fmt.Println(err)
		code := http.StatusInternalServerError
		http.Error(w, http.StatusText(code), code)
		return
	}

	if t, ok := tmpls[reqPath]; ok {
		thisSection := ""
		for _, link := range nav.Links {
			if strings.HasPrefix(reqPath, link.Href) {
				thisSection = link.Href
			}
		}
		nav := nav
		nav.ThisSection = thisSection
		page := Page{nav}

		err := t.Execute(w, page)
		if err != nil {
			fmt.Println(err)
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
			return
		}
	} else {
		http.ServeFile(w, r, fp.Join(root, reqPath))
	}
}

func main() {
	http.HandleFunc("/", rootHandler)
	fmt.Printf("Listening on %s:%s\n", host, port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), nil))
}
