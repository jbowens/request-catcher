package main

import (
	"crypto/subtle"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
	catcher := catcher.NewCatcher(config)

	tlsconf := &tls.Config{MinVersion: tls.VersionTLS10}

	// Start a http server to redirect http traffic to https
	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  30 * time.Second,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			// Catch HTTP requests too, even if we are redirecting them.
			catcher.Catch(req)

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
		Handler:      withPProfHandler(catcher),
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

func withPProfHandler(next http.Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	pprofHandler := basicAuth(mux, os.Getenv("PPROFPW"), "admin")

	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Host == "requestcatcher.com" && strings.HasPrefix(req.URL.Path, "/debug/pprof") {
			pprofHandler.ServeHTTP(rw, req)
			return
		}
		next.ServeHTTP(rw, req)
	})
}

func basicAuth(handler http.Handler, password, realm string) http.Handler {
	p := []byte(password)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, pass, ok := r.BasicAuth()
		if !ok || subtle.ConstantTimeCompare([]byte(pass), p) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			w.WriteHeader(401)
			io.WriteString(w, "Unauthorized\n")
			return
		}
		handler.ServeHTTP(w, r)
	})
}

func fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}
