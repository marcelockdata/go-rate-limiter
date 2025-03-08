package main

import (
	"github.com/marcelockdata/go-rate-limiter/config"
	"github.com/marcelockdata/go-rate-limiter/router"
)

func main() {
	config.Init()
	router.Init()
}
