package main

import (
    "os"
    "io/ioutil"
    "net/http"
    "strings"
    "path/filepath"
    "flag"
    "fmt"
)

var pwd, _ = os.Getwd()
var port = flag.String("p", "8080", "Port to listen. Default 8080")

type FileDetail struct {
	name string
	path string
}

func stat(path string) (isExist bool, info os.FileInfo, _ error) {
    fileInfo, err := os.Stat(path)
    if err == nil { return true, fileInfo, nil }
    if os.IsNotExist(err) { return false, nil, nil }
    return false, nil, err
}

func webPage(path string) (body []byte, status int) {
	exist, info, _ := stat(path)

	if ! exist {
		return  make([]byte, 0), 404
	}

	if ! info.IsDir() {
	   _body, err := ioutil.ReadFile(path)
	   if err != nil { return make([]byte, 0), 500 }
	   return _body, 200
	}

  indexFile := path + string(os.PathSeparator) + "index.html"
  _body, _status := webPage(indexFile)
  if _status == 200 {
    return _body, _status
  }

	files := loadDir(path)
	fileLinks := []string{}

	for _, fileDetail := range files { fileLinks = append(fileLinks, "<a href='/" + fileDetail.path + "'>" + fileDetail.name + "</a>")  }

	return  []byte(strings.Join(fileLinks, " <br/>")), 200
}


func loadDir(path string) ([]FileDetail) {

	_files := []FileDetail{}

    children, _ := ioutil.ReadDir(path)


    for _, f := range children {
    		fullPath := path + string(os.PathSeparator) + f.Name()

    		_files = append(_files, FileDetail{name:  f.Name(), path: fullPath});
    }

	return _files
}

func handler(w http.ResponseWriter, r *http.Request) {

	urlPath := r.URL.Path[1:]

	if urlPath == "" {
		urlPath, _ = filepath.Rel(pwd, pwd)
	}

  body, status := webPage(urlPath)
  w.WriteHeader(status)
	w.Write(body)
}

func main() {
    flag.Parse()
    http.HandleFunc("/", handler)
    fmt.Println("Listing on port :", *port)
    http.ListenAndServe(":" + *port, nil)
}
