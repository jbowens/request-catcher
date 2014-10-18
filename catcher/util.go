package catcher

import "strings"

func hostWithoutPort(host string) string {
	if sepIndex := strings.IndexRune(host, ':'); sepIndex != -1 {
		host = host[:sepIndex]
	}
	return host
}
