package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/danikarik/handler/pkg/service"
)

var addr = flag.String("http.addr", "", "Address for listening")

func main() {
	flag.Parse()

	if *addr == "" {
		*addr = ":" + os.Getenv("PORT")
	}

	var (
		srv  = service.New(*addr)
		errC = make(chan error, 1)
	)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errC <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		log.Println("start listening on: " + *addr)
		errC <- srv.Start()
	}()

	if err := <-errC; err == http.ErrServerClosed {
		log.Println("shutting down server ...")
	} else {
		log.Fatal(err)
	}
}
