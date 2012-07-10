// gomedia is a basic music server written with Go (recently updated to Go 1).
// It is designed to serve music via a web interface from your iTunes library
// but you can point to any directory that contains your music
package main

import (
	"flag"
	"net/http"
	"com.mrd/web"
	"com.mrd/types"
)

const (
	filePrefix = "/f/"
)

var (
	addr    = flag.String("port", ":8080", "http listen address")
	root    = flag.String("root", "/User/<username>/Music/iTunes/", "music root")
	stdRoot = flag.String("user", "osxuser", "your username if your itunes library is in a standard location")
)

func main() {
	flag.Parse()
	web.Setup(*addr, *root, *stdRoot)
	http.HandleFunc("/", web.StaticHandler)
	http.HandleFunc("/static", web.StaticHandler)
	http.Handle(filePrefix, types.AppHandler(web.FileHandler))
	http.ListenAndServe(*addr, nil)
}
