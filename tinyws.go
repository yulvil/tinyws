package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

type mapping struct {
	Directory string
	Path      string
}

func logRequest(r *http.Request) {
	d, _ := httputil.DumpRequestOut(r, false)
	fmt.Printf("%+v\n", d)
}

var portNum int
var directory string
var path string
var https bool

// var cfg string

func init() {
	flag.IntVar(&portNum, "port", 8080, "server port")
	flag.StringVar(&directory, "dir", "./", "directory to serve")
	flag.StringVar(&path, "path", "/", "context path")
	// flag.StringVar(&cfg, "cfg", "", "json configuration")
	flag.BoolVar(&https, "https", false, "run as https (cert.pem and key.pem required)")
	// openssl req -x509 -newkey rsa:2048 -keyout key.pem -out cert.pem -days XXX -nodes
}

func main() {
	flag.Parse()

	fmt.Printf("Serving directory %s on port %d at path %s\n", directory, portNum, path)
	fmt.Printf("http://localhost:%d%s\n", portNum, path)

	http.Handle(path, http.StripPrefix(path, http.FileServer(http.Dir(directory))))
	if https {
		log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", portNum), "cert.pem", "key.pem", nil))
	} else {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", portNum), nil))
	}
}
