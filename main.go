package main

import "github.com/jbowens/request-catcher/catcher"

func main() {
	requestCatcher := catcher.NewCatcher()
	requestCatcher.Start()
}
