package main

// https://stackoverflow.com/questions/35641888/is-it-possible-to-host-multiple-domain-tls-in-golang-with-net-http#35642256

import (
	"crypto/tls"
	"net/http"
	"time"

	"log"
)

func myHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("tls"))
}

func main() {
	t := log.Logger{}
	var err error
	tlsConfig := &tls.Config{}
	tlsConfig.Certificates = make([]tls.Certificate, 3)
	// go http server treats the 0'th key as a default fallback key
	tlsConfig.Certificates[0], err = tls.LoadX509KeyPair("test0.pem", "key.pem")
	if err != nil {
		t.Fatal(err)
	}
	tlsConfig.Certificates[1], err = tls.LoadX509KeyPair("test1.pem", "key.pem")
	if err != nil {
		t.Fatal(err)
	}
	tlsConfig.Certificates[2], err = tls.LoadX509KeyPair("test2.pem", "key.pem")
	if err != nil {
		t.Fatal(err)
	}
	tlsConfig.BuildNameToCertificate()

	http.HandleFunc("/", myHandler)
	server := &http.Server{
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		TLSConfig:      tlsConfig,
	}

	listener, err := tls.Listen("tcp", ":8443", tlsConfig)
	if err != nil {
		t.Fatal(err)
	}
	log.Fatal(server.Serve(listener))
}
