package main

import (
	"github.com/tages/internal/app"
	"github.com/tages/internal/pkg/config"
)

func main() {
	cfg := config.Load()
	app.Run(cfg)
}
