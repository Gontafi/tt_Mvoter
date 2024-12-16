package main

import (
	"context"
	"log"
	"time"
	"tt/config"
	"tt/internal/app"
)

func main() {
	cfg := config.ReadConfig()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := app.RunApp(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
}
