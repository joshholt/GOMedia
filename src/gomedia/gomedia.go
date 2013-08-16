// gomedia is a basic music server written with Go (recently updated to Go 1).
// It is designed to serve music via a web interface from your iTunes library
// but you can point to any directory that contains your music
package main

import (
	"flag"
	"github.com/joshholt/types"
	"github.com/joshholt/web"
	"net/http"
)

// filePrefix defines the path the SPA frontend uses to browse directories and
// load your music files
const (
	filePrefix = "/f/"
)

// Command Line Flags
// addr -> the port on which this server runs
// root -> the optional root location of your music files
// stdRoot -> just provide your username if you store your music in the std iTunes Location
var (
	addr    = flag.String("port", ":8080", "http listen address")
	root    = flag.String("root", "/User/<username>/Music/iTunes/", "music root")
	stdRoot = flag.String("user", "osxuser", "your username if your itunes library is in a standard location")
)

// main kicks everything off by parsing the command line flags, sets up the web
// handlers and starts the httpserver.
func main() {
	flag.Parse()
	web.Setup(*addr, *root, *stdRoot)
	http.HandleFunc("/", web.StaticHandler)
	http.HandleFunc("/static", web.StaticHandler)
	http.Handle(filePrefix, types.AppHandler(web.FileHandler))
	http.ListenAndServe(*addr, nil)
}
