package main

import (
	"log"
	"net/http"
	"os"
)





func main() {
	l := log.New(os.Stdout, "product-api", log.LstdFlags)
	hh := handlers.NewHello(l)

	sm := http.NewServerMux()
	sm.Handle("/", hh)


	http.ListenAndServe(":9090", sm)

}