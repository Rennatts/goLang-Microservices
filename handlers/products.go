package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/renatama/microservices/data"
)

type Products struct {
	l * log.Logger
}

func NewProducts(l * log.Logger) *Products {
	return &Products{l}
}





func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
    p.l.Println("==hereeee====")
    
    lp := data.GetProducts()
    err := lp.ToJSON(rw)
    if err != nil {
        http.Error(rw, "unable to marshal json", http.StatusInternalServerError)
    }
}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
    prod, ok := r.Context().Value(KeyProduct{}).(*data.Product)
    p.l.Println("prod", prod)
    if !ok || prod == nil {
        http.Error(rw, "Invalid product data", http.StatusBadRequest)
        return
    }

    p.l.Println("Adding product:", prod)

    data.AddProduct(prod)
}



func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(rw, "Unable to convert id", http.StatusBadRequest)
        return
    }

    prod := r.Context().Value(KeyProduct{}).(*data.Product)
    if prod == nil {
        http.Error(rw, "Invalid product data", http.StatusBadRequest)
        return
    }

    err = data.UpdateProduct(id, prod)
    if err != nil {
        http.Error(rw, "Product not found", http.StatusNotFound)
        return
    }
}



type KeyProduct struct{}

func (p Products) MiddlewareProductsValidation(next http.Handler) http.Handler {
    return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
        prod := data.Product{}

        // Error handling for product deserialization
        err := prod.FromJSON(r.Body)
        if err != nil {
            p.l.Println("error deserializing", err)
            http.Error(rw, "error reading product", http.StatusBadRequest)
            return
        }

        // Error handling for product validation
        err = prod.Validate()
        if err != nil {
            p.l.Println("error validating product", err)
            http.Error(rw, fmt.Sprintf("error validating product: %v", err), http.StatusBadGateway)
            return
        }

        // Adding the product to the context
        ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
        r = r.WithContext(ctx)

        // Proceed to the next handler
        next.ServeHTTP(rw, r)
    })
}





