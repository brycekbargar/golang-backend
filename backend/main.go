package main

import (
	"github.com/brycekbargar/realworld-backend/ports/echohttp"
)

func main() {
	echohttp.Start(4123, "replace me")
}
