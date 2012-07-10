package types

import (
  "os"
  "time"
  "net/http"
)

type LibraryItem struct {
  Name string `json:"name"`
  Size int64 `json:"size"`
  Mode os.FileMode `json:"mode"`
  ModTime time.Time `json:"modTime"`
  IsDir bool `json:"isDir"`
}

type AppError struct {
  Error error `json:"error"`
  Message string `json:"message"`
  Code int `json:"errorCode"`
}

type AppHandler func(http.ResponseWriter, *http.Request) *AppError

func (fn AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  if err := fn(w, r); err != nil {
    http.Error(w, err.Message, err.Code)
  }
}