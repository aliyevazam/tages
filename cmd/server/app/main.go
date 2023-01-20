package main

import (
	"github.com/aliyevazam/tages/internal/app"
	"github.com/aliyevazam/tages/internal/pkg/config"
)

func main() {
	cfg := config.Load()
	app.Run(cfg)
}
