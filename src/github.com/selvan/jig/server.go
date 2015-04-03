package main

import (
    "os"
    "io/ioutil"
    "net/http"
    "strings"
    "path/filepath"
)

var pwd, _ = os.Getwd()

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

func webPage(path string) ([]byte) {
	exist, info, _ := stat(path)

	if ! exist {
		return  []byte("<html><h3>Error!! Not exist </h3></html>")
	}

	if ! info.IsDir() {
	   body, err := ioutil.ReadFile(path)
	   if err != nil { return []byte("System error :: Unable to processes " + path) }
	   return body
	}

	files := loadDir(path)
	fileLinks := []string{}

	for _, fileDetail := range files { fileLinks = append(fileLinks, "<a href='/" + fileDetail.path + "'>" + fileDetail.name + "</a>")  }

	return  []byte(strings.Join(fileLinks, " <br/>"))
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

	w.Write(webPage(urlPath))
}

func main() {
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}