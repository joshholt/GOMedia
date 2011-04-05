package main

import (
	"flag"
	"http"
	"io"
	"json"
	"os"
	"mime"
	"path"
	"strconv"
	"strings"
)

const (
	filePrefix = "/f/"
)

var (
	addr = flag.String("http", ":8080", "http listen address")
	root = flag.String("root", "/User/<username>/Music/iTunes/", "music root")
)

func main() {
	flag.Parse()
	http.HandleFunc("/", StaticHandler)
	http.HandleFunc("/static", StaticHandler)
	http.HandleFunc(filePrefix, File)
	http.ListenAndServe(*addr, nil)
}

func StaticHandler(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path == "/" {
    http.ServeFile(w, r, "public/index.html")
  } else {
    http.ServeFile(w, r, strings.Replace(r.URL.Path, "/", "", 1)) 
  }
}

func File(w http.ResponseWriter, r *http.Request) {
	fn := *root + r.URL.Path[len(filePrefix):]
	fi, err := os.Stat(fn)
	if err != nil {
		http.Error(w, err.String(), http.StatusNotFound)
		return
	}
	if fi.IsDirectory() {
		serveDirectory(fn, w, r)
		return
	}
	f, err := os.Open(fn, os.O_RDONLY, 0)
	if err != nil {
		http.Error(w, err.String(), http.StatusInternalServerError)
		return
	}
	t := mime.TypeByExtension(path.Ext(fn))
	if t == "" {
		t = "application/octet-stream"
	}
	w.SetHeader("Content-Type", t)
	w.SetHeader("Content-Length", strconv.Itoa64(fi.Size))
	io.Copy(w, f)
}

func serveDirectory(fn string, w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err, ok := recover().(os.Error); ok {
			http.Error(w, err.String(), http.StatusInternalServerError)
		}
	}()
	d, err := os.Open(fn, os.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}
	files, err := d.Readdir(-1)
	if err != nil {
		panic(err)
	}
	j := json.NewEncoder(w)
	w.SetHeader("Content-Type", "application/json")
	if err := j.Encode(files); err != nil {
		panic(err)
	}
}
