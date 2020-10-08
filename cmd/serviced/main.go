package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/danikarik/handler/pkg/service"
)

var addr = flag.String("http.addr", "", "Address for listening")

func main() {
	flag.Parse()

	if *addr == "" {
		*addr = ":" + os.Getenv("PORT")
	}

	var (
		srv = &http.Server{
			Addr:         *addr,
			Handler:      service.New(),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		}
		errC = make(chan error, 1)
	)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errC <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		log.Println("start listening on: " + *addr)
		errC <- srv.ListenAndServe()
	}()

	<-errC

	fmt.Println("")
	log.Println("shutting down server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown failed: %v", err)
	}
}
