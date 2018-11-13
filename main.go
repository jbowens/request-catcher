package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/jbowens/request-catcher/catcher"
)

func main() {
	if len(os.Args) < 2 {
		fatalf("Usage: request-catcher <config-filename>\n")
	}
	config, err := catcher.LoadConfiguration(os.Args[1])
	if err != nil {
		fatalf("error loading configuration file: %s\n", err)
	}

	tlsconf := &tls.Config{MinVersion: tls.VersionTLS10}

	// Start a http server to redirect http traffic to https
	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  30 * time.Second,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Connection", "close")
			url := "https://" + req.Host + req.URL.String()
			http.Redirect(w, req, url, http.StatusMovedPermanently)
		}),
	}
	go func() { log.Fatal(srv.ListenAndServe()) }()

	// Start the HTTPS server.
	fullHost := config.Host + ":" + strconv.Itoa(config.HTTPSPort)
	server := http.Server{
		Addr:         fullHost,
		Handler:      catcher.NewCatcher(config),
		TLSConfig:    tlsconf,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		TLSNextProto: map[string]func(*http.Server, *tls.Conn, http.Handler){},
	}
	// TODO: use GetCertificate instead and periodically reload
	// the tls keypair from disk. as written, certificate renewals
	// require a process restart.
	// NOTE: can't use autocert because we need to use DNS challenges
	// to acquire wildcard certificates.
	err = server.ListenAndServeTLS(
		filepath.Join(config.TLSDir, "fullchain.pem"),
		filepath.Join(config.TLSDir, "privkey.pem"),
	)
	if err != nil {
		fatalf("error listening on %s: %s\n", fullHost, err)
	}
}

func fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}
