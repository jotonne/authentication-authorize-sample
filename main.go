package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	port                    int
	gracefulShutdownTimeout int
)

func main() {
	parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
	})

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           mux,
		ReadHeaderTimeout: 10,
	}
	idleConnectionClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		sig := <-sigint
		fmt.Println(sig.String())
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(gracefulShutdownTimeout))
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("success to shutdown")
		}
		close(idleConnectionClosed)
	}()

	fmt.Printf("start to HTTP Server. Port is %d\n", port)
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("failed to start HTTP Server. Error is %s\n", err.Error())
	}
	<-idleConnectionClosed
}

func parse() {
	flag.IntVar(&port, "port", 8080, "API Server port")
	flag.IntVar(&gracefulShutdownTimeout, "gracefulShutdownTimeout", 30, "GracefulShutdown's Timeout(second)")
	flag.Parse()
}
