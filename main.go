package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var modtime time.Time

func handler(w http.ResponseWriter, r *http.Request) {
	base := filepath.Base(r.URL.Path)

	var filename string
	if base == "/" {
		filename = "index.html"
	} else {
		filename = base
	}

	if filename == "main.wasm.gz" {
		w.Header().Add("Content-Type", "application/wasm")
		w.Header().Add("Content-Encoding", "gzip")
	}

	location := fmt.Sprintf("assets/%s", filename)
	file, err := os.Open(location)
	if err != nil {
		log.Printf("could not find file %s: %v", location, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	http.ServeContent(w, r, filename, modtime, file)
}

func main() {
	host := flag.String("host", ":8000", "ip and port to run server on")
	flag.Parse()

	http.HandleFunc("/", handler)
	log.Printf("running on host %s\n", *host)
	log.Fatal(http.ListenAndServe(*host, nil))
}
