package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
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

	fullHost := config.Host + ":" + strconv.Itoa(config.Port)
	handler := catcher.NewCatcher(config)
	server := http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		TLSNextProto: map[string]func(*http.Server, *tls.Conn, http.Handler){},
		Addr:         fullHost,
		Handler:      handler,
	}
	err = server.ListenAndServe()
	if err != nil {
		fatalf("error listening on %s: %s\n", fullHost, err)
	}
}

func fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}
