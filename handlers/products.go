package handlers

import (
	"encoding/json"
	"log"
	"microservicespetprod/data"
	"net/http"
	"time"
)

type Product struct {
	l *log.Logger
}

func NewProduct (l *log.Logger) *Product {
	return &Product{l}
}

func (p *Product) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Request to / at:", time.Now().UTC())
	lp := data.GetProducts()
	
	d, err := json.Marshal(lp)
	if err != nil {
		http.Error(rw, "Cannot marshal json", http.StatusInternalServerError)
	}

	rw.Write(d)
}