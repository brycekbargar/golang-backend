package main

import (
	"github.com/brycekbargar/realworld-backend/ports/echohttp"
)

func main() {
	server := new(echohttp.Server)
	server.Start(4123)
}
