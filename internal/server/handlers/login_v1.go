package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/dmitrykharchenko95/otus_user/customerrors"
	"github.com/dmitrykharchenko95/otus_user/internal/database/entity"
)

func (h *handler) loginV1() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req LoginV1Req
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf(fmt.Sprintf("decode request body error: %v", err))
			httpError(w, customerrors.ErrDecodeBody)
			return
		}

		var (
			userID int64
			err    error
		)
		if userID, err = h.checkCredentials(r.Context(), req.Email, req.Password); err != nil {
			log.Printf(fmt.Sprintf("check user credentials error: %v", err))
			switch err {
			case entity.ErrInvalidPass:
				httpError(w, customerrors.ErrWrongCredentials)
				return
			case sql.ErrNoRows:
				httpError(w, customerrors.ErrNotFound)
				return
			}
			httpError(w, customerrors.ErrInternal)
			return
		}

		var token string
		if token, err = h.auth.GenerateJWT(userID); err != nil {
			log.Printf(fmt.Sprintf("generate jwt error: %v", err))
			httpError(w, customerrors.ErrInternal)
			return
		}

		if err = json.NewEncoder(w).Encode(
			LoginV1Resp{
				Token:  token,
				UserID: userID,
			},
		); err != nil {
			log.Printf("parse response error: %v", err)
			httpError(w, customerrors.ErrInternal)
			return
		}

		return
	})
}

func (h *handler) checkCredentials(ctx context.Context, email, pass string) (int64, error) {
	var u, err = h.db.GetByEmail(ctx, email)
	if err != nil {
		return 0, err
	}

	return u.Id, u.ValidatePassword(pass)
}
