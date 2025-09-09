package handlers

import (
	"fmt"
	"net/http"
)

func ProductsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Products endpoint not implemented yet.")
}
