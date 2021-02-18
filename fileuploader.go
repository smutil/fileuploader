package main

import (
	"fmt"
	"net/http"
	"os"
	"io"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"flag"
	"log"
	"strings"
)

var workingDir string

func main() {
	port := flag.String("port", "3000", "overwrite default port")
	dest := flag.String("dest", "", "(required) destination directory, should not be root /")
	viewmode := flag.Bool("viewmode", false, "/view will be enabled to view all the files in destination directory")

	flag.Parse()
	if *dest == "" || *dest == "/" {
		log.Fatal("-dest is required for default destination/working directory and should not be root /, please refer -h")
	} 
	workingDir = formatDirName(*dest)
	log.Println("working directory is  "+workingDir)
	fs := http.FileServer(http.Dir(workingDir))
	http.Handle("/metrics", promhttp.Handler())
	if *viewmode {
		log.Println("starting application is view mode")
		http.Handle( "/view/", http.StripPrefix( "/view", fs ) )
	}
	
	http.HandleFunc("/upload", uploadFile)
	http.HandleFunc("/", ping)
	log.Println("application starting on port  "+*port)
	http.ListenAndServe(":"+*port, nil)
	
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	// Maximum upload of 10 MB files
	r.ParseMultipartForm(10 << 20)
	// Get handler for metadata file 
	err := r.ParseForm()

	destDir := r.FormValue("dest")
	if len(destDir) != 0 {
		workingDir += formatDirName(destDir)
	}
	log.Println("file will be uploaded to destination directory "+workingDir)

	// Get handler for filename and size 
	file, handler, err := r.FormFile("data")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	defer file.Close()
	
	log.Printf("Uploaded File: %+v\n", handler.Filename)

	makeDirectoryIfNotExists(workingDir)

	// Create file
	dst, err := os.Create(workingDir+"/"+handler.Filename)
	defer dst.Close()

	if err != nil {
		log.Printf("some error occured while creating the file")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		log.Printf("some error occured while copying the file")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Successfully Uploaded File " + workingDir +"/"+handler.Filename +"\n")
}


func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "/metrics")
	fmt.Fprintln(w, "/upload")
	fmt.Fprintln(w, "/view")
}

func formatDirName(s string) string {
	if strings.HasSuffix(s, "/") { 
		s = strings.TrimSuffix(s, "/")
	}
	if !(strings.HasPrefix(s, "/")){
		s = "/" + s
	}
	return s
}

func makeDirectoryIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil
}
