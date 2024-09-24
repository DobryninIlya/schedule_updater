package main

import (
	"context"
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/jackc/pgx"
	"log"
	"main/internal/app/service"
	"main/internal/app/store"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"
)

var (
	ConfigPath string
)

func init() {
	flag.StringVar(&ConfigPath, "config-path", path.Join("configs", "main.toml"), "path to config file")
}

const gate = "http://127.0.0.1:8383"

type App struct {
	done   chan os.Signal
	Config *store.Config
}

func NewApp() *App {
	config := store.NewConfig()
	_, err := toml.DecodeFile(ConfigPath, config)
	if err != nil {
		log.Fatal(err)
	}
	ret := &App{
		done:   make(chan os.Signal, 1),
		Config: config,
	}
	signal.Notify(ret.done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	return ret
}

func (a *App) run() {
	updater := service.NewUpdater(context.Background(), time.Second*1, pgx.ConnConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     5432,
		Database: os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
	}, a.Config)

	go func() {
		log.Println("Service schedule updater started")
		updater.Run()
	}()
	<-a.done
	log.Println("Exiting")
}

func main() {
	var app = NewApp()
	app.run()

}
