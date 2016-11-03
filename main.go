package main

import (
	"github.com/mholt/caddy/caddy/caddymain"

	// Register Extensions
	_ "github.com/abiosoft/caddy-git"
	_ "github.com/hacdias/caddy-minify"
)

func main() {
	caddymain.Run()
}
