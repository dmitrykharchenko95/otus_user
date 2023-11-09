package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/dmitrykharchenko95/otus_user/customerrors"
	"github.com/dmitrykharchenko95/otus_user/internal/database"
	"github.com/dmitrykharchenko95/otus_user/internal/database/entity"
	"github.com/dmitrykharchenko95/otus_user/internal/server/middlewares"
)

type (
	Handler interface {
		http.Handler
		createV1() http.Handler
		getV1() http.Handler
		deleteV1() http.Handler
		updateV1() http.Handler
	}

	handler struct {
		*mux.Router
		code int
		db   database.Manager
	}
)

func NewHandler(db database.Manager) http.Handler {
	var h = handler{
		Router: mux.NewRouter(),
		db:     db,
	}
	h.Router.Handle("/user", h.createV1()).Methods(http.MethodPost)
	h.Router.Handle("/user/{id}", h.getV1()).Methods(http.MethodGet)
	h.Router.Handle("/user/{id}", h.updateV1()).Methods(http.MethodPut)
	h.Router.Handle("/user/{id}", h.deleteV1()).Methods(http.MethodDelete)

	h.Router.Use(middlewares.Default...)

	h.Router.Handle("/metrics", promhttp.Handler())
	h.Router.Handle("/health", h.healthCheck())
	return h
}

func (h *handler) createV1() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req CreateV1Req
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf(fmt.Sprintf("decode request body error: %v", err))
			httpError(w, customerrors.ErrDecodeBody)
			return
		}

		var id, err = h.db.Add(r.Context(), &entity.User{
			Username:  req.Username,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
			Phone:     req.Phone,
		})
		if err != nil {
			log.Printf(fmt.Sprintf("create user in database error: %v", err))
			httpError(w, customerrors.ErrInternal)
			return
		}

		log.Printf("User created - id %d\n", id)

		if err = json.NewEncoder(w).Encode(map[string]int64{"id": id}); err != nil {
			log.Printf(fmt.Sprintf("parse response error: %v", err))
			httpError(w, customerrors.ErrInternal)
			return
		}

	})
}

func (h *handler) getV1() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var id, err = strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			log.Printf(fmt.Sprintf("parsed user id error: %v", err))
			httpError(w, customerrors.ErrParseQuery)
			return
		}

		var u *entity.User
		if u, err = h.db.Get(r.Context(), int64(id)); err != nil {
			log.Printf(fmt.Sprintf("get user from database error: %v", err))
			if errors.Is(err, sql.ErrNoRows) {
				httpError(w, customerrors.ErrNotFound)
				return
			}
			httpError(w, customerrors.ErrInternal)
			return
		}

		if err = json.NewEncoder(w).Encode(
			GetV1Resp{
				Id:        u.Id,
				Username:  u.Username,
				FirstName: u.FirstName,
				LastName:  u.LastName,
				Email:     u.Email,
				Phone:     u.Phone,
			},
		); err != nil {
			log.Printf(fmt.Sprintf("parse response error: %v", err))
			httpError(w, customerrors.ErrInternal)
			return
		}
	})
}

func (h *handler) deleteV1() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var id, err = strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			log.Printf(fmt.Sprintf("parsed user id error: %v", err))
			httpError(w, customerrors.ErrParseQuery)
			return
		}

		if err = h.db.Delete(r.Context(), int64(id)); err != nil {
			log.Printf(fmt.Sprintf("delete user from database error: %v", err))
			httpError(w, customerrors.ErrInternal)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}

func (h *handler) updateV1() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req UpdateV1Req
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf(fmt.Sprintf("decode request body error: %v", err))
			httpError(w, customerrors.ErrDecodeBody)
			return
		}

		var id, err = strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			log.Printf(fmt.Sprintf("parsed user id error: %v", err))
			httpError(w, customerrors.ErrParseQuery)
			return
		}

		if _, err = h.db.Get(r.Context(), int64(id)); err != nil {
			log.Printf(fmt.Sprintf("get user from database error: %v", err))
			if errors.Is(err, sql.ErrNoRows) {
				httpError(w, customerrors.ErrNotFound)
				return
			}
			httpError(w, customerrors.ErrInternal)
			return
		}

		if err = h.db.Update(r.Context(), &entity.User{
			Id:        int64(id),
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
			Phone:     req.Phone,
		}); err != nil {
			log.Printf(fmt.Sprintf("add user to database error: %v", err))
			httpError(w, customerrors.ErrInternal)
			return
		}
	})
}

func (h *handler) healthCheck() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := h.db.Ping(); err != nil {
			log.Printf(fmt.Sprintf("ping database error: %v", err))
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
