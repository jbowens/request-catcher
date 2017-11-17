package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

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
					fmt.Fprintf(os.Stderr, "Error listening on %s: %s\n", fullHost, err)
					os.Exit(1)
				}
			},
		},
	}
	cliApp.Run(os.Args)
}
