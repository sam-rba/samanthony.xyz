package main

import (
	"flag"
	"fmt"
	"golang.org/x/sys/unix"
	"html/template"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	fp "path/filepath"
	"strings"
)

// Flags
var (
	host   = "localhost"
	port   = "80"
	chroot = "/var/www/"
	user   = "www"
	group  = "www"
	root   = "/htdocs/samanthony.xyz/"
)

func init() {
	flag.StringVar(&host, "host", host, "")
	flag.StringVar(&port, "port", port, "")
	flag.StringVar(&chroot, "chroot", chroot, "")
	flag.StringVar(&user, "user", user, "")
	flag.StringVar(&group, "group", group, "")
	flag.StringVar(&root, "root", root, "")

	flag.Parse()
}

// Must lookup the hostname before entering the chroot.
var addr = ""

func init() {
	// host is an ip address
	if ip := net.ParseIP(host); ip != nil {
		addr = ip.String()
	} else { // host is a domain name
		addrs, err := net.LookupHost(host)
		if err != nil {
			log.Fatal(err)
		}
		for _, a := range addrs {
			if ip := net.ParseIP(a); ip != nil {
				if v4 := ip.To4(); v4 != nil {
					addr = v4.String()
				}
			}
		}
		if addr == "" {
			log.Fatalf("No ipv4 address bound to %s", host)
		}
	}
}

var (
	uid int
	gid int
)

func init() {
	var err error
	uid, err = uidOf(user)
	if err != nil {
		log.Fatal(err)
	}
	gid, err = gidOf(group)
	if err != nil {
		log.Fatal(err)
	}
}

// Enter chroot
func init() {
	if err := unix.Chroot(chroot); err != nil {
		log.Fatalf("chroot: %s: %v", chroot, err)
	}
}

// Build templates
var tmpl = make(map[string]*template.Template)

func init() {
	err := fp.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if fp.Clean(path) == fp.Clean(root) ||
			fp.Ext(path) != ".html" ||
			path == fp.Join(root, "base.html") {
			return nil
		}
		label := path[len(fp.Clean(root)):]
		tmpl[label] = template.Must(template.ParseFiles(fp.Join(root, "base.html"), path))
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

// Template data
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
	if err := dropPerms(uid, gid); err != nil {
		log.Println(err)
		code := http.StatusInternalServerError
		http.Error(w, http.StatusText(code), code)
		return
	}

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
		log.Println(err)
		code := http.StatusInternalServerError
		http.Error(w, http.StatusText(code), code)
		return
	}

	if t, ok := tmpl[reqPath]; ok {
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
		http.ServeFile(w, r, fp.Join(root, reqPath))
	}
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.Handle("/.well-known/acme-challenge/",
		http.StripPrefix(
			"/.well-known/acme-challenge/",
			http.FileServer(http.Dir("/acme/")),
		),
	)

	log.Printf("Listening on %s:%s\n", addr, port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", addr, port), nil))
}
