package main

import (
	"github.com/MHSaeedkia/tinyTimer/internal/router"
)

const (
	addr = "localhost"
	port = "19191"
)

func main() {
	router.GetEngine().Run(addr + ":" + port)

	var forever chan struct{}
	<-forever
}
