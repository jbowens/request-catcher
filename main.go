package main

import (
	"fmt"
	"os"
	"strconv"

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
				bindHost := args[0]
				port, _ := strconv.Atoi(args[1])
				rootHost := args[2]
				fmt.Println(bindHost)
				requestCatcher := catcher.NewCatcher(bindHost, port, rootHost)
				requestCatcher.Start()
			},
		},
	}
	cliApp.Run(os.Args)
}
