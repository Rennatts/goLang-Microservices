package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/renatama/microservices/data"
)

type Products struct {
	l * log.Logger
}

func NewProducts(l * log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("aaaaaaaaaaa")
    if r.Method != http.MethodGet {
        http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    lp := data.GetProducts()
    d, err := json.Marshal(lp)
    if err != nil {
        p.l.Println("[ERROR] marshaling json:", err)
        http.Error(rw, "unable to marshal json", http.StatusInternalServerError)
        return
    }

    rw.Header().Set("Content-Type", "application/json")
    rw.Write(d)
}
