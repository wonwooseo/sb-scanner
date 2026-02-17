package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"

	pkglog "sb-scanner/pkg/logger"
	"sb-scanner/router"
)

func main() {
	cfgF := flag.String("config", "", "path to config file")
	portF := flag.Int("port", 80, "port to listen (default: 80)")
	debugF := flag.Bool("debug", false, "enable debug mode")
	flag.Parse()

	// read in config
	v := viper.New()
	v.SetConfigFile(*cfgF)
	if err := v.ReadInConfig(); err != nil {
		slog.Error("failed to read config", "err", err)
		os.Exit(1)
	}

	// set up logger
	pkglog.InitLogger(v.GetString("loglevel"))

	// server initialization
	logger := pkglog.GetLogger().With("pkg", "main")
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", *portF),
		Handler:      router.New(v, *debugF),
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}
	serverCtx, serverCancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		logger.Info("got signal")

		shutdownCtx, shutdownCancel := context.WithTimeout(serverCtx, 10*time.Second)
		go func() {
			defer shutdownCancel()
			<-shutdownCtx.Done() // if serverCtx is not cancelled before shutdownCtx times out, forcefully exit
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				logger.Error("graceful shutdown timeout; exiting forcefully")
				os.Exit(1)
			}
		}()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Error("graceful shutdown failed; exiting forcefully", "err", err)
		}
		serverCancel()
	}()

	logger.Info("start listening..", "addr", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		logger.Error("stopping server", "err", err)
	}

	<-serverCtx.Done() // wait for graceful shutdown
	logger.Info("exiting..")
}
