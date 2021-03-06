package main

import (
	"fmt"
	"net/http"
	"os"
	"io"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"flag"
	"log"
	"strings"
)

var workingDir, authToken string
var Version = "v1.5"

var (
	endpointsAccessed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "endpoints_accessed",
			Help: "Total number of accessed to a given endpoint",
		},
		[]string{"accessed_endpoint"},
	)
)

func main() {
	version := flag.Bool("version", false, "returns the fileuploader version")
	port := flag.String("port", "3000", "overwrite default port")
	dest := flag.String("dest", "", "(required) destination directory, should not be root /")
	viewmode := flag.Bool("viewmode", false, "/view will be enabled to view all the files in destination directory")
	tlsCert := flag.String("tls-crt", "", "certificate path, only needed for ssl service")
	tlsKey := flag.String("tls-key", "", "key path, only needed for ssl service")
	logfile := flag.String("log-file", "", "key path, only needed for ssl service")
	authtoken := flag.String("authToken", "", "Bearer Toekn for basic authentication")
	prometheus.MustRegister(endpointsAccessed)
	flag.Parse()
	if *version {
		fmt.Println(Version)
		return
	}
	if *dest == "" || *dest == "/" {
		log.Fatal("-dest is required for default destination/working directory and should not be root /, please refer -h")
	} 
	if *authtoken != "" {
		authToken = *authtoken
	}

	if *logfile != "" {
		file, err := os.OpenFile(*logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(file)
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
	http.HandleFunc("/", health)
	http.HandleFunc("/health", health)
	
	if *tlsCert != "" && *tlsKey != "" {
		log.Println("application starting on port  "+*port + " (https)")
		err := http.ListenAndServeTLS(":"+*port, *tlsCert, *tlsKey, nil)
		if err  != nil {
			log.Println(err)
            return
		}
	} else {
		log.Println("application starting on port  "+*port + " (http)")
		err := http.ListenAndServe(":"+*port, nil)
		if err  != nil {
			log.Println(err)
        	return
		}
	}
}

func authenticateRequest(w http.ResponseWriter, req *http.Request) bool {
	var authenticated bool = false
	if authToken != "" {
		authorizationHeader := req.Header.Get("authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				if bearerToken[1] == authToken {
					authenticated = true
				} else {
					http.Error(w, "Invalid authorization token", http.StatusUnauthorized)
					return authenticated
				}
			} else {
				http.Error(w, "Invalid authorization token", http.StatusUnauthorized)
				return authenticated
			}
		} else {
			http.Error(w, "An authorization header not provided", http.StatusUnauthorized)
			return authenticated
		}
	} else {
		authenticated = true		
	}
	return authenticated
}

func view(w http.ResponseWriter, r *http.Request) {
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	endpointsAccessed.WithLabelValues("/uploadFile").Inc()
	var finalDestDir string
	if !authenticateRequest(w,r) {
		return
	}
	
	// Maximum upload of 10 MB files
	r.ParseMultipartForm(10 << 20)
	// Get handler for metadata file 
	err := r.ParseForm()
	destDir := r.FormValue("dest")
	if len(destDir) != 0 {
		finalDestDir += workingDir + formatDirName(destDir)
	} else {
		finalDestDir = workingDir
	} 
	// Get handler for filename and size 
	file, handler, err := r.FormFile("data")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	defer file.Close()
	log.Printf("received File: %+v\n", handler.Filename)
	log.Println("file will be uploaded to destination directory "+finalDestDir)


	makeDirectoryIfNotExists(finalDestDir)

	// Create file
	dst, err := os.Create(finalDestDir+"/"+handler.Filename)
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
	log.Println("Successfully Uploaded File : " + finalDestDir +"/"+handler.Filename +"\n")
	fmt.Fprintf(w, "Successfully Uploaded File : " + handler.Filename +"\n")
}

func health(w http.ResponseWriter, r *http.Request) {
	endpointsAccessed.WithLabelValues("/health").Inc()
	fmt.Fprintln(w, "application is healthy")
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
