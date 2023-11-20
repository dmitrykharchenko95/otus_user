package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dmitrykharchenko95/otus_user/customerrors"
	"github.com/dmitrykharchenko95/otus_user/internal/database/entity"
)

func (h *handler) createV1() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req CreateV1Req
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("decode request body error: %v", err)
			httpError(w, customerrors.ErrDecodeBody)
			return
		}

		var u = &entity.User{
			Username:  req.Username,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
			Phone:     req.Phone,
		}

		if len(strings.TrimSpace(req.Password)) < 6 {
			log.Println("password is too short")
			httpError(w, customerrors.ErrShortPassword)
			return
		}

		if err := u.SetPassword(req.Password); err != nil {
			log.Printf("set password error: %v", err)
			httpError(w, customerrors.ErrInternal)
			return
		}

		var id, err = h.db.Add(r.Context(), u)
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
