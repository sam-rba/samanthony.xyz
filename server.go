package main

import (
	"flag"
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

var (
	host   = "samanthony.xyz"
	port   = "443"
	htdocs = "/var/www/htdocs/samanthony.xyz"
)

const (
	acmeDocs = "/var/www/acme/"
	certFile = "/etc/ssl/samanthony.xyz.fullchain.pem"
	keyFile  = "/etc/ssl/private/samanthony.xyz.key"
)

const (
	devHost   = "localhost"
	devPort   = "6969"
	devHtdocs = "htdocs/"
)

var devMode bool

func init() {
	flag.BoolVar(&devMode, "dev", false,
		"Run server in debug/development mode (on localhost without tls)")

	flag.Parse()

	if devMode {
		host = devHost
		port = devPort
		htdocs = devHtdocs
	}
}

var tmpls = make(map[string]*tmpl.Template)

func init() {
	err := fp.WalkDir(htdocs, func(path string, d fs.DirEntry, err error) error {
		if fp.Clean(path) == fp.Clean(htdocs) ||
			fp.Ext(path) != ".html" ||
			path == fp.Join(htdocs, "base.html") {
			return nil
		}
		label := path[len(fp.Clean(htdocs)):]
		tmpls[label] = tmpl.Must(tmpl.ParseFiles(fp.Join(htdocs, "base.html"), path))
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
	if info, err := os.Stat(fp.Join(htdocs, reqPath)); err == nil {
		if info.IsDir() {
			reqPath = path.Join(reqPath, "index.html")
		}
	} else if os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	} else {
		log.Println(err)
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
			log.Println(err)
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
			return
		}
	} else {
		http.ServeFile(w, r, fp.Join(htdocs, reqPath))
	}
}

func main() {
	http.HandleFunc("/", rootHandler)
	if !devMode {
		http.Handle("/.well-known/acme-challenge/",
			http.StripPrefix(
				"/.well-known/acme-challenge/",
				http.FileServer(http.Dir(acmeDocs)),
			),
		)
	}

	if devMode {
		log.Printf("Listening on %s:%s\n", devHost, devPort)
		log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", devHost, devPort), nil))
	} else {
		log.Printf("Listening on %s:%s\n", host, port)
		log.Fatal(http.ListenAndServeTLS(fmt.Sprintf("%s:%s", host, port), certFile, keyFile, nil))
	}
}
