package handlers

import (
	"context"
	"encoding/json"
	"io/ioutil"
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


// func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
//     product := &data.Product{}

//     p.l.Println("======", product, "=========")

//     err := product.FromJSON(r.Body)
//     if err != nil {
//         http.Error(rw, "unable to unmarshal json", http.StatusBadRequest)
//     }

//     data.AddProduct(product)
// }

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
    prod, ok := r.Context().Value(KeyProduct{}).(*data.Product)
    if !ok || prod == nil {
        http.Error(rw, "Invalid product data", http.StatusBadRequest)
        return
    }

    p.l.Println("Adding product:", prod)

    // Proceed with adding the product
    data.AddProduct(prod)

    // Respond with success or added product information
}




func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r);

    id, err := strconv.Atoi(vars["id"]);
    if err != nil {
        http.Error(rw, "unable to convert id", http.StatusBadRequest)
        return
    }


    // Fetch the existing product
    prod := r.Context().Value(KeyProduct{}).(*data.Product)


    p.l.Println("prod", prod)

    err = data.UpdateProduct(id, prod)
    if err != nil {
        http.Error(rw, "Product not found", http.StatusNotFound)
        return
    }

    if err != nil {
        http.Error(rw, "Product not found", http.StatusNotFound)
        return
    }
}


type KeyProduct struct{}


func (p Products) MiddlewareProductsValidation(next http.Handler) http.Handler {
    return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
        // Read and parse the body here only once
        body, err := ioutil.ReadAll(r.Body)
        if err != nil {
            http.Error(rw, "Error reading request body", http.StatusBadRequest)
            return
        }

        prod := &data.Product{}
        err = json.Unmarshal(body, prod)
        if err != nil {
            http.Error(rw, "unable to unmarshal json", http.StatusBadRequest)
            return 
        }

        ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
        req := r.WithContext(ctx)

        next.ServeHTTP(rw, req)
    })
}




