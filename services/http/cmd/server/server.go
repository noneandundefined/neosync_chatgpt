package main

import (
	"context"
	"fmt"
	"log"
	"neomatica/neosync/config"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

var (
	Version = "1.0.x"
)

func (s *httpServer) httpStart() error {
	port, err := strconv.Atoi(config.HttpServerPort)
	if err != nil {
		log.Fatalf("Invalid port: %v", err)
	}

	/* Initial HTTPx routes */
	routes := s.routes()

	fmt.Printf("\n[%v] [INFO] HTTP server started :%d\n", time.Now().Format("2006-01-02 15:04:05"), port)
	fmt.Printf("[%v] [INFO] Proccess PID: %d, Version: %s\n", time.Now().Format("2006-01-02 15:04:05"), os.Getpid(), Version)
	fmt.Printf("[%v] [INFO] Golang version: %s\n\n", time.Now().Format("2006-01-02 15:04:05"), runtime.Version())

	httpServe := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           routes,
		ReadTimeout:       5 * time.Minute,
		WriteTimeout:      5 * time.Minute,
		IdleTimeout:       90 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	idleConnsClosed := make(chan struct{})

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c

		fmt.Printf("\n[%v] [INFO] HTTP server shutting down...\n", time.Now().Format("2006-01-02 15:04:05"))

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err := httpServe.Shutdown(ctx); err != nil {
			fmt.Printf("\n[%v] [ERROR] HTTP shutdown error: %s\n", time.Now().Format("2006-01-02 15:04:05"), err.Error())
		}

		close(idleConnsClosed)
	}()

	if err := httpServe.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	<-idleConnsClosed
	return nil
}
