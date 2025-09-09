package handlers

import (
	"fmt"
	"net/http"
)

func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Categories endpoint not implemented yet.")
}
