package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/acme/autocert"

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

	handler := catcher.NewCatcher(config)
	srv := http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		TLSNextProto: map[string]func(*http.Server, *tls.Conn, http.Handler){},
		Addr:         config.ListenAddress(),
		Handler:      handler,
	}

	// If there's no Let's Encrypt configuration, just launch
	// without TLS.
	if config.LetsEncrypt == nil {
		err = srv.ListenAndServe()
		if err != nil {
			fatalf("error listening on %s: %s\n", config.ListenAddress(), err)
		}
		return
	}

	// Setup an autocert configuration to dynamically acquire
	// certificates for requested domains.
	m := autocert.Manager{
		Prompt:      autocert.AcceptTOS,
		Cache:       autocert.DirCache(config.LetsEncrypt.CertsDirectory),
		HostPolicy:  subdomainsHostPolicy(config.RootHost),
		RenewBefore: 24 * time.Hour,
	}
	srv.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}

	err = srv.ListenAndServeTLS("", "")
	if err != nil {
		fatalf("error listening on %s: %s\n", config.ListenAddress(), err)
	}
}

func fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func subdomainsHostPolicy(rootHost string) autocert.HostPolicy {
	return func(_ context.Context, host string) error {
		if host == rootHost {
			return nil
		}
		if strings.HasSuffix(host, "."+rootHost) {
			return nil
		}
		return errors.New("host doesn't match root host")
	}
}
