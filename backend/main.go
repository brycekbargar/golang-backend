package main

import (
	"github.com/brycekbargar/realworld-backend/adapters/inmemory"
	"github.com/brycekbargar/realworld-backend/ports"
	"github.com/brycekbargar/realworld-backend/ports/echohttp"
)

func main() {
	// TODO: Configure the port and secret
	im := inmemory.NewInstance()
	echohttp.Start(
		ports.DefaultJWTConfig("Replace Me"),
		4123,
		im.Users,
		im.Articles,
	)
}
