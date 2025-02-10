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

	"github.com/gin-gonic/gin"
)

var (
	port                    int
	gracefulShutdownTimeout int
)

func main() {
	parse()

	router := gin.Default()
	router.GET("/health", func(ctx *gin.Context) {
		ctx.Header("Content-Type", "application/json; charset=utf-8")
		ctx.String(http.StatusOK, "ok")
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
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
