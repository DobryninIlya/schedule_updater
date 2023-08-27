package main

import (
	"context"
	"github.com/jackc/pgx"
	"log"
	"main/internal/service"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const gate = "http://127.0.0.1:8383"

type App struct {
	done chan os.Signal
}

func NewApp() *App {
	ret := &App{
		done: make(chan os.Signal, 1),
	}
	signal.Notify(ret.done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	return ret
}

func (a *App) run() {
	//db.MakeConnection()
	updater := service.NewUpdater(context.Background(), time.Minute*1, pgx.ConnConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     5432,
		Database: os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
	})

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
