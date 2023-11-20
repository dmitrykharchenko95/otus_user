package handlers

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"github.com/dmitrykharchenko95/otus_user/customerrors"
)

func (h *handler) deleteV1() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var token = r.Header.Get("Authorization")
		if token == "" {
			log.Printf("header Authorization not set")
			httpError(w, customerrors.ErrUnauthorized)
			return
		}

		var ok bool
		if token, ok = strings.CutPrefix(token, "Bearer "); !ok {
			log.Printf("header Authorization has not Bearer prefix")
			httpError(w, customerrors.ErrWrongToken)
			return
		}

		var claims, err = h.auth.CheckToken(token)
		if err != nil {
			log.Printf("check token err: %v", err)
			httpError(w, customerrors.ErrWrongToken)
			return
		}

		var id int
		if id, err = strconv.Atoi(mux.Vars(r)["id"]); err != nil {
			log.Printf("parsed user id error: %v", err)
			httpError(w, customerrors.ErrParseQuery)
			return
		}

		if claims.UserID != int64(id) {
			log.Printf("userIDs mismatch")
			httpError(w, customerrors.ErrWrongToken)
			return
		}

		if err = h.db.Delete(r.Context(), int64(id)); err != nil {
			log.Printf("delete user from database error: %v", err)
			httpError(w, customerrors.ErrInternal)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}
