// Package rest provides REST HTTP requests handlers.
package rest

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

	"book-management-system/docs"
	"book-management-system/entities/constants"
	"book-management-system/usecases"
)

// Init REST controllers
func Init(useCase *usecases.UseCase) {
	r := mux.NewRouter()

	NewBookController(r, useCase)

	initDoc(r)
	serve(r)
}

func initDoc(r *mux.Router) {
	docs.SwaggerInfo.Title = constants.ServiceName
	docs.SwaggerInfo.Version = constants.ServiceVersion
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
}

func serve(r http.Handler) {
	var wait time.Duration
	flag.DurationVar(
		&wait,
		"graceful-timeout",
		time.Second*15,
		"the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m",
	)
	flag.Parse()

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Error HTTP server shutdown: %v", err)
	}

	log.Println("Shutting down server gracefully")
	os.Exit(0)
}
