package main

import (
  "encoding/json"
  "flag"
  "io"
  "mime"
  "net/http"
  "os"
  "time"
  "path"
  "strconv"
  "strings"
)

const (
  filePrefix = "/f/"
)

type LibraryItem struct {
  Name string `json:"name"`
  Size int64 `json:"size"`
  Mode os.FileMode `json:"mode"`
  ModTime time.Time `json:"modTime"`
  IsDir bool `json:"isDir"`
}

var (
  addr = flag.String("port", ":8080", "http listen address")
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
    http.Error(w, err.Error(), http.StatusNotFound)
    return
  }
  if fi.IsDir() {
    serveDirectory(fn, w, r)
    return
  }
  f, err := os.Open(fn)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  t := mime.TypeByExtension(path.Ext(fn))
  if t == "" {
    t = "application/octet-stream"
  }
  w.Header().Set("Content-Type", t)
  w.Header().Set("Content-Length", strconv.FormatInt(fi.Size(), 10))
  io.Copy(w, f)
}

func handleNonHTTPError(err error) {
  if err != nil {
    panic(err)
  }
}

func findFiles(fn string) []os.FileInfo {
  dir, err := os.Open(fn)
  handleNonHTTPError(err)
  items, err := dir.Readdir(-1)
  handleNonHTTPError(err)
  return items
}

func buildLibrary(items []os.FileInfo) []LibraryItem {
  var library = []LibraryItem{}

  for i := 0; i < len(items); i++ {
    f := items[i]
    library = append(library, LibraryItem{f.Name(), f.Size(), f.Mode(), f.ModTime(), f.IsDir()})
  }

  return library
}

func serveDirectory(fn string, w http.ResponseWriter, r *http.Request) {
  defer func() {
    if err, ok := recover().(error); ok {
      http.Error(w, err.Error(), http.StatusInternalServerError)
    }
  }()

  var library = buildLibrary(findFiles(fn))

  j := json.NewEncoder(w)
  w.Header().Set("Content-Type", "application/json")
  if err := j.Encode(library); err != nil {
    panic(err)
  }
}
