package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sb-scanner/route"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.HideBanner = true

	route.AddRoutes(e)

	// graceful exit logics
	errc := make(chan error)     // channel to receive error caused by server
	sigc := make(chan os.Signal) // channel to receive signal
	signal.Notify(sigc, syscall.SIGTERM, syscall.SIGINT)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	go func() {
		if err := e.Start(":80"); err != nil {
			errc <- err
		}
	}()

	select {
	case <-sigc:
		log.Println("gracefully exiting..")
		if err := e.Shutdown(ctx); err != nil {
			e.Logger.Fatal(err)
		}
		log.Println("bye!")
	case err := <-errc:
		log.Fatal(err)
	}
}
