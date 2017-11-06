package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {

	http.HandleFunc("/upload", upload)
	http.Handle("/download/", http.StripPrefix("/download/", http.FileServer(http.Dir("./"))))

	log.Fatal(http.ListenAndServe(":8081", nil))
}

func upload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	lfile, handler, err := r.FormFile("uploadfile")
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	defer lfile.Close()

	filename := handler.Filename

	// mkdir path
	if strings.Contains(filename, "/") {
		lastSep := strings.LastIndex(filename, "/")
		if lastSep > 0 {
			err = os.MkdirAll(filename[:lastSep], 0777)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
		}
	}

	f, err := os.Create(filename)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	defer f.Close()

	io.Copy(f, lfile)
	fmt.Fprintln(w, "upload ok!")
}
