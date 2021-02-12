package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/phaesoo/shield/apps/shield"
	"github.com/phaesoo/shield/configs"
)

type App interface {
	Listen() error
	Shutdown() error
}

func runApp(app App, onDone func()) {
	done := make(chan struct{})
	shutdown := make(chan os.Signal, 1)

	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-shutdown

		if err := app.Shutdown(); err != nil {
		}
		close(done)
	}()

	if err := app.Listen(); err != http.ErrServerClosed {
		panic(err)
	}

	<-done
	onDone()
}

func main() {
	log.Print("Run")

	wg := sync.WaitGroup{}

	app := shield.NewApp(configs.Get())
	wg.Add(1)
	go runApp(app, wg.Done)

	log.Print("Wait")
	wg.Wait()
	log.Print("Finished")
}
