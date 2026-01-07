package handlers

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
)

type VerifyPaymentRequest struct {
	RazorpayOrderID   string `json:"razorpay_order_id"`
	RazorpayPaymentID string `json:"razorpay_payment_id"`
	RazorpaySignature string `json:"razorpay_signature"`
}

func VerifyPaymentHandler(w http.ResponseWriter, r *http.Request) {
	var req VerifyPaymentRequest
	json.NewDecoder(r.Body).Decode(&req)

	secret := os.Getenv("RAZORPAY_KEY_SECRET")

	data := req.RazorpayOrderID + "|" + req.RazorpayPaymentID

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	if expectedSignature != req.RazorpaySignature {
		http.Error(w, "Invalid payment signature", http.StatusBadRequest)
		return
	}

	_, _ = Db.Exec(context.Background(),
		`UPDATE orders
		 SET razorpay_payment_id=$1, payment_status='PAID'
		 WHERE razorpay_order_id=$2`,
		req.RazorpayPaymentID, req.RazorpayOrderID,
	)

	json.NewEncoder(w).Encode(map[string]string{
		"status": "payment_verified",
	})
}
