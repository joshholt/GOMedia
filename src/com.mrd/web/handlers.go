// GoMedia's web package contains all of the logic and HTTPHandlers need
// to serve the music from your iTunes library or any other place on your
// computer's HD
package web

import (
	"com.mrd/types"
	"net/http"
	"encoding/json"
	"io"
	"mime"
	"os"
	"path"
	"strconv"
	"strings"
)

const (
	filePrefix = "/f/"
)

var (
	addr    *string
	root    *string
	stdRoot *string
)

// Setup should be the first function executed after calling flag.parse() in
// the main function of the gomedia program. It receives the parsed command
// line flags, for use inside the handlers
func Setup(port string, fullPath string, user string) {
	addr = &port
	root = &fullPath
	stdRoot = &user
}

func buildLibrary(fn string) (library []types.LibraryItem, err error) {

	dir, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	items, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(items); i++ {
		f := items[i]
		library = append(library, types.LibraryItem{f.Name(), f.Size(), f.Mode(), f.ModTime(), f.IsDir()})
	}

	return library, err
}

// StaticHandler is just what it says it is. It serves the static (public) files
// that are used to display GoMedia's UI/Logic
func StaticHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.ServeFile(w, r, "public/index.html")
	} else {
		http.ServeFile(w, r, strings.Replace(r.URL.Path, "/", "", 1))
	}
}

// FileHandler is the HTTPHandler that determines if you have selected a directory
// or a music file, that reacts accordingly. If you have selected a directory
// then FileHandler calls the serveDirectory function to reveal it's contents.
// If you have selected a music file, FileHandler will copy it's contents into
// the http.ResponseWriter's buffer (i.e. writer)
func FileHandler(w http.ResponseWriter, r *http.Request) *types.AppError {
	if *stdRoot != "osxuser" {
		*root = "/Users/" + *stdRoot + "/Music/iTunes/iTunes Media/Music/"
	}
	fn := *root + r.URL.Path[len(filePrefix):]
	fi, err := os.Stat(fn)

	if err != nil {
		return &types.AppError{err, "File Not Found", http.StatusNotFound}
	}

	if fi.IsDir() {
		return serveDirectory(fn, w, r)
	}

	f, err := os.Open(fn)
	if err != nil {
		return &types.AppError{err, "Could Not Open the file", http.StatusInternalServerError}
	}
	defer f.Close()

	t := mime.TypeByExtension(path.Ext(fn))
	if t == "" {
		t = "application/octet-stream"
	}

	w.Header().Set("Content-Type", t)
	w.Header().Set("Content-Length", strconv.FormatInt(fi.Size(), 10))
	io.Copy(w, f)
	return nil
}

func serveDirectory(fn string, w http.ResponseWriter, r *http.Request) *types.AppError {
	encoder := json.NewEncoder(w)
	library, err := buildLibrary(fn)

	if err != nil {
		return &types.AppError{err, "Could not find files", http.StatusInternalServerError}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := encoder.Encode(library); err != nil {
		return &types.AppError{err, "Could Not encode the JSON", http.StatusInternalServerError}
	}
	return nil
}
