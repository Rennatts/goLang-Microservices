// This line declares the package name as handlers. In Go, each file belongs to a package, 
// which is a way to group functions, and it's the entry point to access the functions, types, etc., 
// defined in the file.
package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// This defines a new struct type named Hello with a single field l, which is a pointer to a log.Logger. 
// This log.Logger will be used for logging.
type Hello struct {
	l *log.Logger
}


// This function named NewHello is a constructor for the Hello type. It takes a pointer to a log.
// Logger as an argument and returns a new instance of Hello with the logger.
func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

// This method makes Hello an implementer of the http.Handler interface. The method has two parameters: 
// rw (an http.ResponseWriter) for writing the response, and r (an *http.Request) for the incoming HTTP request.

func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.l.Println("hello world")

	// These lines read all data from the request body (r.Body). If an error occurs 
	// (e.g., if the body is not readable), it sends an HTTP 400 Bad Request error response and returns from 
	// the function.
	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "opps", http.StatusBadRequest)
		return 
	}

	fmt.Fprintf(rw, "hello %s", d)
}