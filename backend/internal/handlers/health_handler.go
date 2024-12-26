package handlers

import (
	"fmt"
	"net/http"
)

func ChechHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "It's alive!\n")
}
