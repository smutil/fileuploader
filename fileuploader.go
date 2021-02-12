package main

import (
	"fmt"
	"net/http"
	"os"
	"io"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"flag"
	"log"

)

var destDir string
func main() {
	port := flag.String("port", "3000", "overwrite default port")
	dest := flag.String("dest", "", "(required) destination directory")
	
	flag.Parse()
	if *dest == "" {
		log.Fatal("-dest is required for destination directory, please refer -h")
	} 
	destDir = *dest
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", ping)
	http.HandleFunc("/upload", uploadFile)
	http.ListenAndServe(":"+*port, nil)
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	// Maximum upload of 10 MB files
	r.ParseMultipartForm(10 << 20)
	// Get handler for filename, size and headers
	file, handler, err := r.FormFile("data")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size (in bytes): %+v\n", handler.Size)

	// Create file
	dst, err := os.Create(destDir+"/"+handler.Filename)
	defer dst.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Successfully Uploaded File\n")
}

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "/metrics")
	fmt.Fprintln(w, "/upload")
}
