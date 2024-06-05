package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"go-lang-final/internal/models"

	"go-lang-final/internal/store"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func RegisterRESTHandlers(r *mux.Router, store *store.PaymentStore, logger *logrus.Logger) {
	r.HandleFunc("/payments", createPaymentHandler(store, logger)).Methods("POST")
	r.HandleFunc("/payments/{id}", getPaymentHandler(store, logger)).Methods("GET")
	r.HandleFunc("/payments/{id}", updatePaymentHandler(store, logger)).Methods("PUT")
	r.HandleFunc("/payments/{id}", deletePaymentHandler(store, logger)).Methods("DELETE")
	r.HandleFunc("/payments", listPaymentsHandler(store, logger)).Methods("GET")
}

func createPaymentHandler(store *store.PaymentStore, logger *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payment models.Payment
		if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
			logger.Error("Failed to decode request body:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := store.CreatePayment(payment); err != nil {
			logger.Error("Failed to create payment:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(payment)
	}
}

func getPaymentHandler(store *store.PaymentStore, logger *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := strconv.Atoi(vars["id"])
		payment, err := store.GetPayment(id)
		if err != nil {
			logger.Error("Failed to get payment:", err)
			http.Error(w, "Payment not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(payment)
	}
}

func updatePaymentHandler(store *store.PaymentStore, logger *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := strconv.Atoi(vars["id"])
		var payment models.Payment
		if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
			logger.Error("Failed to decode request body:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := store.UpdatePayment(id, payment); err != nil {
			logger.Error("Failed to update payment:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(payment)
	}
}

func deletePaymentHandler(store *store.PaymentStore, logger *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := strconv.Atoi(vars["id"])
		if err := store.DeletePayment(id); err != nil {
			logger.Error("Failed to delete payment:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func listPaymentsHandler(store *store.PaymentStore, logger *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filter := r.URL.Query().Get("filter")
		sort := r.URL.Query().Get("sort")
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

		payments, err := store.ListPayments(filter, sort, page, pageSize)
		if err != nil {
			logger.Error("Failed to list payments:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(payments)
	}
}
