package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/jbowens/request-catcher/catcher"
)

const defaultPort = 80

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

				requestCatcher := catcher.NewCatcher(config)
				err = requestCatcher.Start()
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
			},
		},
	}
	cliApp.Run(os.Args)
}
