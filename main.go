// Package main is the entrypoint for the server.
package main

import "github.com/eltorocorp/badgeserver/badgeserver"

const (
	port = "80"
)

func main() {
	badgeserver.Run(port)
}
