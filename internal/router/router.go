package router

import (
	"go-lang-final/internal/handlers"
	"go-lang-final/internal/store"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func NewRouter(store *store.PaymentStore, logger *logrus.Logger) *mux.Router {
	r := mux.NewRouter()
	handlers.RegisterRESTHandlers(r, store, logger)
	return r
}
