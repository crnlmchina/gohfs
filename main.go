// Go http file server
//
// Author: wangyuxuan
// health check url: /
package main

import (
	"net/http"
	"os"
	"log"
	"fmt"
	"io/ioutil"
	"strings"
)

var baseStorageLocation string

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Base storage location is required.")
	}
	baseStorageLocation = os.Args[1]

	http.HandleFunc("/upload", handleUpload)
	http.HandleFunc("/", handleHealthCheck)

	log.Fatal(http.ListenAndServe(":9090", nil))
}

func handleHealthCheck(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprint(writer, "ok")
}

func handleUpload(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
	} else {
		err := request.ParseForm()

		if err != nil {
			http.Error(writer, "", http.StatusInternalServerError)
			return
		}

		file := request.Form["file"][0]

		body := request.Body
		data, err := ioutil.ReadAll(body)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}

		filePath := baseStorageLocation + file
		folder := filePath[:strings.LastIndex(filePath, "/")]

		if _, error := os.Stat(folder); os.IsNotExist(error) {
			os.MkdirAll(folder, 0755)
		}

		err = ioutil.WriteFile(baseStorageLocation + file, data, 0644)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	}
}