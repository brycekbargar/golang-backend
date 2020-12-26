package main

import (
	"github.com/brycekbargar/realworld-backend/ports"
	"github.com/brycekbargar/realworld-backend/ports/echohttp"
)

func main() {
	echohttp.Start(ports.DefaultJWTConfig("Replace Me"), 4123, nil)
}
