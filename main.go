package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"yak/internal/config"

	"github.com/gorilla/mux"
)

const (
	configPath = "./config.yaml"
)

func main() {

	config, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	router := mux.NewRouter()
	s := &http.Server{
		Addr:    ":" + fmt.Sprintf("%v", config.Server.Port),
		Handler: router,
	}
	fmt.Printf("http server started on %d ...", config.Server.Port)

	//gracefull shutdown
	connClose := make(chan struct{})
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sig

		err = s.Shutdown(context.Background())
		if err != nil {
			log.Panic(err.Error())
		}
		close(connClose)
	}()

	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("server shut down error %v", err)
	}
	<-connClose
	fmt.Println("server stopped...")
}
