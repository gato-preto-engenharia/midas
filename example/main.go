package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/gato-preto-engenharia/midas"
)

type App struct{}

func (App) Run(_ context.Context, cfg midas.Config) error {
	slog.Info(fmt.Sprintf("hello %s!", cfg.Get("name", "World")))

	return nil
}

func main() {
	config := midas.NewConfigFromJson([]byte(`{"name": "Midas"}`))

	midas.SetupSlog(config)
	midas.Run(context.Background(), config, &App{})
}
