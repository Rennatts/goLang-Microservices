package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/renatama/microservices/data"
)

type Products struct {
	l * log.Logger
}

func NewProducts(l * log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
    p.l.Printf("Request received: Method: %s, Path: %s\n", r.Method, r.URL.Path)

    switch r.Method {
    case http.MethodGet:
        p.getProducts(rw, r)
    case http.MethodPost:
        p.addProduct(rw, r)
    case http.MethodPut:
        p.updateProduct(rw, r)
    default:
        http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
    }
}



func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
    lp := data.GetProducts()
    err := lp.ToJSON(rw)
    if err != nil {
        http.Error(rw, "unable to marshal json", http.StatusInternalServerError)
    }
}


func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
    product := &data.Product{}

    p.l.Println("======", product, "=========")

    err := product.FromJSON(r.Body)
    if err != nil {
        http.Error(rw, "unable to unmarshal json", http.StatusBadRequest)
    }

    data.AddProduct(product)
}


// func extractIDFromURL(r *http.Request, l *log.Logger) (int, error) {
//     parts := strings.Split(r.URL.Path, "/")
//     l.Printf("URL Parts: %#v\n", parts)

//     if len(parts) < 3 {
//         return 0, fmt.Errorf("invalid URL path")
//     }

//     id, err := strconv.Atoi(parts[2])
//     if err != nil {
//         l.Printf("Error converting ID: %v\n", err)
//         return 0, fmt.Errorf("invalid ID")
//     }

//     return id, nil
// }

func extractIDFromURL(r *http.Request) (int, error) {
    parts := strings.Split(r.URL.Path, "/")
    if len(parts) < 3 {
        return 0, fmt.Errorf("invalid URL path")
    }
    return strconv.Atoi(parts[2]) // Assuming the ID is the third part of the URL
}


// func (p *Products) updateProduct(rw http.ResponseWriter, r *http.Request) {
//     p.l.Println("hereeee")
//     // Extract the product ID from the URL
//     id, err := extractIDFromURL(r)
//     p.l.Println("here", id)
//     if err != nil {
//         http.Error(rw, "Invalid URL", http.StatusBadRequest)
//         return
//     }

//     // Deserialize the incoming JSON payload into a Product struct
//     prod := &data.Product{}
//     err = prod.FromJSON(r.Body)
//     if err != nil {
//         http.Error(rw, "unable to unmarshal json", http.StatusBadRequest)
//         return
//     }

//     // Update the product
//     data.UpdateProduct(id, prod)
// }

func (p *Products) updateProduct(rw http.ResponseWriter, r *http.Request) {
    p.l.Println("Update product request received")

    // Extract the product ID from the URL
    id, err := extractIDFromURL(r)
    if err != nil {
        http.Error(rw, "Invalid URL", http.StatusBadRequest)
        return
    }

    // Fetch the existing product
    existingProd, err := data.GetProductByID(id)
    if err != nil {
        http.Error(rw, "Product not found", http.StatusNotFound)
        return
    }

    // Deserialize the incoming JSON payload into a temporary Product struct
    updatedProd := &data.Product{}
    err = updatedProd.FromJSON(r.Body)
    if err != nil {
        http.Error(rw, "unable to unmarshal json", http.StatusBadRequest)
        return
    }

    // Apply changes to the existing product
    if updatedProd.Name != "" {
        existingProd.Name = updatedProd.Name
    }
    // Repeat for other fields that can be updated

    // Update the product in your data store
    err = data.UpdateProduct(id, existingProd)
    if err != nil {
        http.Error(rw, "Error updating product", http.StatusInternalServerError)
        return
    }
}


