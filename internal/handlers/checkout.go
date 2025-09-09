package handlers

import (
	"fmt"
	"net/http"
)

func CheckoutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Checkout endpoint not implemented yet.")
}
