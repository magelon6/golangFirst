package main

import (
	"context"
	"log"
	"microservicespetprod/handlers"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/swaggo/http-swagger"
	_ "microservicespetprod/docs"
)

//	@title			Coffe Shop API
//	@version		1.0
//	@description	This is my 100 attempt to add swaggeer to this proj :)
//	@host			localhost:9090
//	@BasePath		/v1
func main() {
	l := log.New(os.Stdout, "product-api ", log.LstdFlags)
	ph := handlers.NewProduct(l)
	uh := handlers.NewUser(l)

	sm := mux.NewRouter()

	// Product handlers

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", ph.GetProducts)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProducts)
	putRouter.Use(ph.MiddlewareValidateProduct)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", ph.AddProduct)
	postRouter.Use(ph.MiddlewareValidateProduct)

	// User handlers

	getUserRouter := sm.Methods(http.MethodGet).Subrouter()
	getUserRouter.HandleFunc("/users", uh.GetUsers)

	postUserRouter := sm.Methods(http.MethodPost).Subrouter()
	postUserRouter.HandleFunc("/users", uh.AddUser)
	postUserRouter.Use(uh.ValidateMiddlewareUser)

	putUserRouter := sm.Methods(http.MethodPost).Subrouter()
	putUserRouter.HandleFunc("/users/{id:[0-9]+}", uh.UpdateUser)
	putUserRouter.Use(uh.ValidateMiddlewareUser)

	// swagger router

	sm.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
        httpSwagger.URL("/swagger/doc.json"), 
    ))

	s := http.Server{
		Addr: "localhost:9090",
		Handler: sm,
		IdleTimeout: 120 * time.Second,
		WriteTimeout: 1*time.Second,
		ReadTimeout: 1*time.Second,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan

	l.Println("Recieved terminate, graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
