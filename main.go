package main

import "jackson/webhook-catcher/catcher"

func main() {
	webhooksCatcher := catcher.NewCatcher()
	webhooksCatcher.Start()
}
