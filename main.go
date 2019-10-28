package main

import (
	"github.com/caddyserver/caddy/caddy/caddymain"

	// Register Extensions
	_ "github.com/abiosoft/caddy-git"
	_ "github.com/hacdias/caddy-minify"
)

func main() {
	caddymain.EnableTelemetry = false
	caddymain.Run()
}
