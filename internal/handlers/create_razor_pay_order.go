package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
)

type RazorpayOrderResponse struct {
	ID       string `json:"id"`
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}

func CreateRazorpayOrderHandler(w http.ResponseWriter, r *http.Request) {
	userIDVal := r.Context().Value("user_id")
	if userIDVal == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := userIDVal.(int)

	var req CheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	// 🔹 calculate total (same logic as checkout)
	var total float64
	for _, item := range req.Items {
		var price float64
		err := Db.QueryRow(context.Background(),
			"SELECT price FROM products WHERE id=$1",
			item.ProductID,
		).Scan(&price)
		if err != nil {
			http.Error(w, "Invalid product", http.StatusBadRequest)
			return
		}
		total += price * float64(item.Quantity)
	}

	amountPaise := int(total * 100)

	payload := map[string]interface{}{
		"amount":   amountPaise,
		"currency": "INR",
	}

	body, _ := json.Marshal(payload)

	reqRazor, _ := http.NewRequest(
		"POST",
		"https://api.razorpay.com/v1/orders",
		bytes.NewBuffer(body),
	)
	reqRazor.SetBasicAuth(
		os.Getenv("RAZORPAY_KEY_ID"),
		os.Getenv("RAZORPAY_KEY_SECRET"),
	)
	reqRazor.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(reqRazor)
	if err != nil {
		http.Error(w, "Razorpay error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var razorResp RazorpayOrderResponse
	json.NewDecoder(resp.Body).Decode(&razorResp)

	// 🔹 store order as CREATED
	_, _ = Db.Exec(context.Background(),
		`INSERT INTO orders (user_id, total_amount, razorpay_order_id, payment_status)
		 VALUES ($1, $2, $3, 'CREATED')`,
		userID, total, razorResp.ID,
	)

	json.NewEncoder(w).Encode(razorResp)
}
