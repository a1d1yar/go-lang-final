package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"go-lang-final/internal/models"
	"go-lang-final/internal/store"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func RegisterRESTHandlers(r *mux.Router, store *store.PaymentStore, logger *logrus.Logger) {
	handler := NewRestHandler(store)
	r.HandleFunc("/create", handler.CreatePayment).Methods("POST")
	r.HandleFunc("/get", handler.GetPayment).Methods("GET")
	r.HandleFunc("/update", handler.UpdatePayment).Methods("PUT")
	r.HandleFunc("/delete", handler.DeletePayment).Methods("DELETE")
	r.HandleFunc("/list", handler.ListPayments).Methods("GET")
}

type RestHandler struct {
	store *store.PaymentStore
}

func NewRestHandler(store *store.PaymentStore) *RestHandler {
	return &RestHandler{store: store}
}

func (h *RestHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var payment models.Payment
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.store.CreatePayment(payment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *RestHandler) GetPayment(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64) // Convert to int64
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	payment, err := h.store.GetPayment(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(payment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *RestHandler) UpdatePayment(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64) // Convert to int64
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var payment models.Payment
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.store.UpdatePayment(id, payment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *RestHandler) DeletePayment(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64) // Convert to int64
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.store.DeletePayment(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *RestHandler) ListPayments(w http.ResponseWriter, r *http.Request) {
	// Handle query parameters for pagination and filters
	currency := r.URL.Query().Get("currency")
	amountStr := r.URL.Query().Get("amount")
	amount, err := strconv.ParseFloat(amountStr, 64) // Convert to float64
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	pageSizeStr := r.URL.Query().Get("pageSize")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	payments, err := h.store.ListPayments(currency, fmt.Sprintf("%.2f", amount), page, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(payments); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
