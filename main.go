package main

import (
	"log"
	"microservicespetprod/handlers"
	"net/http"
	"os"
)

func main() {
	l := log.New(os.Stdout, "product-api \n", log.LstdFlags)
	hh := handlers.NewHello(l)
	gh := handlers.NewGoodbye(l)

	sm := http.NewServeMux()
	sm.Handle("/", hh)
	sm.Handle("/goodbye", gh)

	http.ListenAndServe(":9090", sm)
}
