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

	"github.com/codegangsta/cli"
	"github.com/jbowens/request-catcher/catcher"
)

func main() {
	cliApp := cli.NewApp()
	cliApp.Name = "Request Catcher"
	cliApp.Usage = "The Request Catcher web server"

	cliApp.Commands = []cli.Command{
		{
			Name:      "start",
			ShortName: "s",
			Usage:     "Start the Request Catcher web server",
			Action: func(c *cli.Context) {
				args := c.Args()
				if len(args) < 1 {
					fmt.Println("Must provide config filename")
					os.Exit(1)
				}
				configFilename := args.First()

				config, err := catcher.LoadConfiguration(configFilename)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
					os.Exit(1)
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
						fmt.Fprintf(os.Stderr, "Error listening on %s: %s\n", config.ListenAddress(), err)
						os.Exit(1)
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
					fmt.Fprintf(os.Stderr, "Error listening on %s: %s\n", config.ListenAddress(), err)
					os.Exit(1)
				}
			},
		},
	}
	cliApp.Run(os.Args)
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
