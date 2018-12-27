package main

import (
	"github.com/mholt/caddy/caddy/caddymain"

	// Register Extensions
	_ "github.com/abiosoft/caddy-git"
	_ "github.com/hacdias/caddy-minify"
	_ "github.com/zikes/gopkg"
)

func main() {
	caddymain.EnableTelemetry = false
	caddymain.Run()
}
