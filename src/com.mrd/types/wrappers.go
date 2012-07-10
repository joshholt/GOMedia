// The types and wrapper types used in gomedia
package types

import (
	"os"
	"time"
	"net/http"
)

// LibraryItem is a struct representation of FileInfo for one file in your
// file system
type LibraryItem struct {
	Name    string      `json:"name"`
	Size    int64       `json:"size"`
	Mode    os.FileMode `json:"mode"`
	ModTime time.Time   `json:"modTime"`
	IsDir   bool        `json:"isDir"`
}

// AppError is a convience type used to return HTTP status codes and messages
// from AppHandler wrapped HTTPHandlers
type AppError struct {
	Error   error  `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"errorCode"`
}

// AppHandler is a fun type that is used to enable generic error handling
// for HTTP handlers
type AppHandler func(http.ResponseWriter, *http.Request) *AppError

// ServerHTTP is the way to wrap your HTTP Handlers in generic error handling
// you just need to return and AppError or nil
func (fn AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		http.Error(w, err.Message, err.Code)
	}
}
