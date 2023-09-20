package handlers

import (
	"context"
	"fmt"
	"log"
	"microservicespetprod/db"

	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Product struct {
	l *log.Logger
}

func NewProduct(l *log.Logger) *Product {
	return &Product{l}
}

// ShowProducts godoc
//
//	@Summary		Show a products list
//	@Description	Get all products
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	db.Product
//
//	@failure		400	{string}	"string"	"error"
//
//	@Router			/ [get]
func (p *Product) GetProducts(rw http.ResponseWriter, r *http.Request) {

	p.l.Println("Request time on / is: ", time.Now().UTC())
	lp := db.GetProducts()
	err := lp.ToJSON(rw)

	if err != nil {
		http.Error(rw, "Unable to marshall JSON", http.StatusInternalServerError)
	}
}

func (p *Product) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Products")

	prod := r.Context().Value(KeyProduct{}).(db.Product)
	db.AddProduct(&prod)
}

func (p Product) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}
	p.l.Println("Handle PUT Product", id)
	prod := r.Context().Value(KeyProduct{}).(db.Product)

	err = db.UpdateProduct(id, &prod)

	if err == db.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}

type KeyProduct struct{}

func (p Product) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := db.Product{}

		err := prod.FromJSON(r.Body)
		if err != nil {
			p.l.Println("[ERROR] deserializing product", err)
			http.Error(rw, "Error reading product", http.StatusBadRequest)
			return
		}
		err = prod.Validate()
		if err != nil {
			p.l.Println("[ERROR] validating product", err)
			http.Error(
				rw,
				fmt.Sprintf("Error validating product: %s", err),
				http.StatusBadRequest,
			)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)

		next.ServeHTTP(rw, req)
	})
}
