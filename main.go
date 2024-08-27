package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

var Proxy *httputil.ReverseProxy

func main() {
	var (
		tunnel string
		port   int
	)
	flag.StringVar(&tunnel, "url", "", "http tunnel url")
	flag.IntVar(&port, "port", 9090, "listening port")
	flag.Parse()

	targetURL, err := url.Parse(tunnel)
	if err != nil {
		panic(err)
	}

	Proxy = httputil.NewSingleHostReverseProxy(targetURL)
	Proxy.Director = func(req *http.Request) {
		req.URL.Scheme = targetURL.Scheme
		req.URL.Host = targetURL.Host
		req.URL.Path = targetURL.Path
	}

	http.HandleFunc("/", proxy)
	http.HandleFunc("/sql", sql)
	log.Println("Listening on port", port)
	log.Println("Proxying to", tunnel)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		panic(err)
	}
}

func proxy(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	boundary := strings.TrimPrefix(r.Header.Get("Content-Type"), "multipart/form-data; boundary=")
	r.Body = io.NopCloser(io.MultiReader(bytes.NewReader(body), r.Body))

	data := parseFormData(body, boundary)

	switch data.Action {
	case "C":
		connect(w, r, data)
	case "Q":
		query(w, r, data)
	}
}

type writer struct {
	http.ResponseWriter
	data []byte
}

func (w *writer) Write(data []byte) (int, error) {
	w.data = append(w.data, data...)
	return w.ResponseWriter.Write(data)
}
