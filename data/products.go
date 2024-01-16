package data

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// Product defines the structure for an API product
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	SKU         string  `json:"sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

// Products is a collection of Product
type Products []*Product

// ToJSON serializes the contents of the collection to JSON
// NewEncoder provides better performance than json.Unmarshal as it does not
// have to buffer the output into an in memory slice of bytes
// this reduces allocations and the overheads of the service
//
// https://golang.org/pkg/encoding/json/#NewEncoder
func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Product) FromJSON(r io.Reader) error {
    e := json.NewDecoder(r)
    return e.Decode(p)
}

// GetProducts returns a list of products
func GetProducts() Products {
	return productList
}

func AddProduct(p *Product) {
    p.ID = getNextId()
    productList = append(productList, p)
}

func UpdateProduct(id int, updatedProd *Product) error {
    for i, prod := range productList {
        if prod.ID == id {
            // Only update fields that are non-zero in the updated product
            if updatedProd.Name != "" {
                productList[i].Name = updatedProd.Name
            }
            if updatedProd.Description != "" {
                productList[i].Description = updatedProd.Description
            }
            if updatedProd.Price != 0 {
                productList[i].Price = updatedProd.Price
            }
            if updatedProd.SKU != "" {
                productList[i].SKU = updatedProd.SKU
            }

            // Update the UpdatedOn field
            productList[i].UpdatedOn = time.Now().UTC().String()

            return nil
        }
    }
    return fmt.Errorf("Product with ID %d not found", id)
}


func getNextId() int {
    lp := productList[len(productList) - 1]
    lp.ID++
    return lp.ID
}

func GetProductByID(id int) (*Product, error) {
    for _, prod := range productList {
        if prod.ID == id {
            return prod, nil
        }
    }
    return nil, fmt.Errorf("Product with ID %d not found", id)
}


// productList is a hard coded list of products for this
// example data source
var productList = []*Product{
    {
        ID:          1,
        Name:        "Latte",
        Description: "Frothy milky coffee",
        Price:       2.45,
        SKU:         "abc323",
        CreatedOn:   time.Now().UTC().String(),
        UpdatedOn:   time.Now().UTC().String(),
    },
    {
        ID:          2,
        Name:        "Espresso",
        Description: "Short and strong coffee without milk",
        Price:       1.99,
        SKU:         "fjd34",
        CreatedOn:   time.Now().UTC().String(),
        UpdatedOn:   time.Now().UTC().String(),
    },
}


