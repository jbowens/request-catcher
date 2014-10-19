package main

import (
	"os"
	"strconv"

	"github.com/jbowens/request-catcher/catcher"
)

const defaultPort = 80

func main() {
	args := os.Args
	host := args[1]
	port, err := strconv.Atoi(args[2])
	if err != nil {
		port = defaultPort
	}
	rootHost := args[3]

	requestCatcher := catcher.NewCatcher(host, port, rootHost)
	requestCatcher.Start()
}
