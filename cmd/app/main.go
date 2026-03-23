package main

import (
	"fmt"

	"github.com/I-Van-Radkov/vesta-gkh/internal/app"
	"github.com/I-Van-Radkov/vesta-gkh/internal/config"
)

func main() {
	cfg, err := config.ParseConfigFromEnv()
	if err != nil {
		panic(fmt.Errorf("failed to parse config: %w", err))
	}

	app.RunApp(cfg)
}
