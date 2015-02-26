package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"

	"github.com/bradfitz/http2"
)

type Config struct {
	HttpPort          int       `json:"http_port,omitempty"`
	HttpsPort         int       `json:"https_port,omitempty"`
	SslCertificate    string    `json:"ssl_certificate,omitempty"`
	SslCertificateKey string    `json:"ssl_certificate_key,omitempty"`
	Handlers          []Handler `json:"handlers,omitempty"`
}

type Handler struct {
	Type        string `json:"type,omitempty"` // file, proxy
	ContextPath string `json:"context_path,omitempty"`
	Directory   string `json:"directory,omitempty"`
	Backend     string `json:"backend,omitempty"`
}

func timeTrack(start time.Time, r *http.Request) {
	elapsed := time.Since(start)
	ms := float64(elapsed) / float64(time.Millisecond)
	fmt.Printf("%s %s %s %s %v\n", r.RemoteAddr, r.Method, r.URL, r.Proto, ms)
}

func WebLog(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer timeTrack(time.Now(), r)
		handler.ServeHTTP(w, r)
	})
}

//func addCORS(handler http.Handler) http.Handler {
//  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//      w.Header().Set("Access-Control-Allow-Origin", "*")
//      w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With")
//      handler.ServeHTTP(w, r)
//  })
//}

var cfg string

func init() {
	flag.StringVar(&cfg, "cfg", "{}", "json configuration")
	// openssl req -x509 -newkey rsa:2048 -keyout key.pem -out cert.pem -days XXX -nodes
}

func loadConfig(cfg string) *Config {
	config := new(Config)

	if cfg[0] == '{' {
		err := json.NewDecoder(bytes.NewBuffer([]byte(cfg))).Decode(config)
		fmt.Printf("%+v\n", err)
	} else {
		buf, err := ioutil.ReadFile(cfg)
		if err != nil {
			log.Fatal(err)
		}
		err = json.NewDecoder(bytes.NewBuffer(buf)).Decode(config)
		fmt.Printf("%+v\n", err)
	}
	fmt.Printf("%+v\n", config)

	return config
}

func main() {
	flag.Parse()

	config := loadConfig(cfg)

	wg := &sync.WaitGroup{}

	for _, handler := range config.Handlers {
		if handler.Type == "file" {
			http.Handle(handler.ContextPath, WebLog(http.StripPrefix(handler.ContextPath, http.FileServer(http.Dir(handler.Directory)))))
		} else if handler.Type == "proxy" {
			u, _ := url.Parse(handler.Backend)
			proxy := httputil.NewSingleHostReverseProxy(u)
			http.Handle(handler.ContextPath, WebLog(proxy))
		}
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.HttpPort), nil))
	}()

	go func() {
		defer wg.Done()
		var srv http.Server
		srv.Addr = fmt.Sprintf(":%d", config.HttpsPort)
		http2.ConfigureServer(&srv, &http2.Server{})
		log.Fatal(srv.ListenAndServeTLS("cert.pem", "key.pem"))
		//log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", config.HttpsPort), "cert.pem", "key.pem", nil))
	}()

	wg.Wait()
}
