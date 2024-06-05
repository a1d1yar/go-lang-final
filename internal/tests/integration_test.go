package tests

import (
	"bytes"
	"encoding/json"
	"go-lang-final/internal/models"
	"go-lang-final/internal/router"
	"go-lang-final/internal/store"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func setupTestServer() *httptest.Server {
	logger := logrus.New()
	dsn := "postgresql://postgres:Aldiyar2004@localhost:5432/library"
	paymentStore, err := store.NewPaymentStore(dsn)
	if err != nil {
		logger.Fatalf("Failed to connect to the database: %v", err)
	}

	r := router.NewRouter(paymentStore, logger)
	return httptest.NewServer(r)
}

func TestCreateAndGetPayment(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	payment := models.Payment{ID: 1, Amount: 100.0, Currency: "USD"}
	body, _ := json.Marshal(payment)
	resp, err := http.Post(server.URL+"/create", "application/json", bytes.NewBuffer(body))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	resp, err = http.Get(server.URL + "/get?id=1")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var gotPayment models.Payment
	json.NewDecoder(resp.Body).Decode(&gotPayment)
	assert.Equal(t, payment, gotPayment)
}

func TestCreateAndUpdatePayment(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	payment := models.Payment{ID: 2, Amount: 150.0, Currency: "EUR"}
	body, _ := json.Marshal(payment)
	resp, err := http.Post(server.URL+"/create", "application/json", bytes.NewBuffer(body))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	updatedPayment := models.Payment{ID: 2, Amount: 200.0, Currency: "EUR"}
	body, _ = json.Marshal(updatedPayment)
	req, _ := http.NewRequest(http.MethodPut, server.URL+"/update?id=2", bytes.NewBuffer(body))
	client := &http.Client{}
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = http.Get(server.URL + "/get?id=2")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var gotPayment models.Payment
	json.NewDecoder(resp.Body).Decode(&gotPayment)
	assert.Equal(t, updatedPayment, gotPayment)
}

func TestCreateAndDeletePayment(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	payment := models.Payment{ID: 3, Amount: 300.0, Currency: "GBP"}
	body, _ := json.Marshal(payment)
	resp, err := http.Post(server.URL+"/create", "application/json", bytes.NewBuffer(body))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	req, _ := http.NewRequest(http.MethodDelete, server.URL+"/delete?id=3", nil)
	client := &http.Client{}
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = http.Get(server.URL + "/get?id=3")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestListPayments(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	payment1 := models.Payment{ID: 4, Amount: 400.0, Currency: "USD"}
	body, _ := json.Marshal(payment1)
	resp, err := http.Post(server.URL+"/create", "application/json", bytes.NewBuffer(body))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	payment2 := models.Payment{ID: 5, Amount: 500.0, Currency: "USD"}
	body, _ = json.Marshal(payment2)
	resp, err = http.Post(server.URL+"/create", "application/json", bytes.NewBuffer(body))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	resp, err = http.Get(server.URL + "/list?currency=USD&amount=400.00&page=1&pageSize=10")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var payments []models.Payment
	json.NewDecoder(resp.Body).Decode(&payments)
	assert.Len(t, payments, 2)
	assert.Equal(t, []models.Payment{payment1, payment2}, payments)
}

func TestCreateAndListPayments(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	payment := models.Payment{ID: 6, Amount: 600.0, Currency: "JPY"}
	body, _ := json.Marshal(payment)
	resp, err := http.Post(server.URL+"/create", "application/json", bytes.NewBuffer(body))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	resp, err = http.Get(server.URL + "/list?currency=JPY&amount=600.00&page=1&pageSize=10")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var payments []models.Payment
	json.NewDecoder(resp.Body).Decode(&payments)
	assert.Len(t, payments, 1)
	assert.Equal(t, payment, payments[0])
}
