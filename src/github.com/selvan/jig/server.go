package main

import (
    "fmt"
    "os"
    "io/ioutil"
    "net/http"
    "strings"
)

func stat(path string) (isExist bool, info os.FileInfo, _ error) {
    fileInfo, err := os.Stat(path)
    if err == nil { return true, fileInfo, nil }
    if os.IsNotExist(err) { return false, nil, nil }
    return false, nil, err
}

func webPage(path string) (string) {
	exist, info, _ := stat(path)

	if ! exist {
		return "<html><h3>Error!! Not exist </h3></html>"
	}

	if ! info.IsDir() {
	   body, err := ioutil.ReadFile(path)
	   if err != nil { return "System error :: Unable to processes " + path }
	   return string(body)
	}

	files := loadDir(path)

	for index, filename := range files { files[index] = "<a href='" + filename + "'>" + filename+ "</a>" }

	return strings.Join(files, " <br/>")
}


func loadDir(path string) ([]string) {	

	_files := []string{}

    files, _ := ioutil.ReadDir(path)
    for _, f := range files {
    		_files = append(_files, f.Name());            
    }

	return _files
}

func handler(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path[1:]
	if urlPath == "" {
		urlPath = "./"
	}

	fmt.Fprintf(w, webPage(urlPath))
}

func main() {
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}