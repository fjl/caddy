package main

import (
	caddycmd "github.com/caddyserver/caddy/v2/cmd"

	// Register extensions
	_ "github.com/abiosoft/caddy-exec"
	_ "github.com/abiosoft/caddy-hmac"
	_ "github.com/caddyserver/caddy/v2/modules/standard"
)

func main() {
	caddycmd.Main()
}
