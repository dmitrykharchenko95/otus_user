package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/dmitrykharchenko95/otus_user/customerrors"
	"github.com/dmitrykharchenko95/otus_user/internal/auth"
	"github.com/dmitrykharchenko95/otus_user/internal/database"
	"github.com/dmitrykharchenko95/otus_user/internal/server/middlewares"
)

type (
	Handler interface {
		http.Handler
		createV1() http.Handler
		loginV1() http.Handler
		getV1() http.Handler
		deleteV1() http.Handler
		updateV1() http.Handler
	}

	handler struct {
		*mux.Router
		code int
		db   database.Manager
		auth *auth.Manager
	}
)

func NewHandler(db database.Manager, jwtKey string) http.Handler {
	var h = handler{
		Router: mux.NewRouter(),
		db:     db,
		auth:   auth.NewManager(jwtKey),
	}
	h.Router.Handle("/login", h.loginV1()).Methods(http.MethodPost)
	h.Router.Handle("/user", h.createV1()).Methods(http.MethodPost)
	h.Router.Handle("/user/{id}", h.getV1()).Methods(http.MethodGet)
	h.Router.Handle("/user/{id}", h.updateV1()).Methods(http.MethodPut)
	h.Router.Handle("/user/{id}", h.deleteV1()).Methods(http.MethodDelete)

	h.Router.Use(middlewares.Default...)

	h.Router.Handle("/metrics", promhttp.Handler())
	h.Router.Handle("/health", h.healthCheck())
	return h
}

func (h *handler) healthCheck() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := h.db.Ping(); err != nil {
			log.Printf("ping database error: %v", err)
			httpError(w, customerrors.New(err, http.StatusInternalServerError))
		}
	})
}

func httpError(w http.ResponseWriter, err customerrors.Error) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(err.Code)
	w.Write(err.GetJSON())
}
