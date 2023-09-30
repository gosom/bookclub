package api

import (
	"encoding/json"
	"net/http"

	"github.com/gosom/bookclub"
	"github.com/gosom/bookclub/api/schema"
)

func (b *bookClubAPI) PostAuthLogin(w http.ResponseWriter, r *http.Request) {
	var payload schema.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		renderError(r, w, bookclub.ErrInvalidBody, "")

		return
	}

	access, refresh, err := b.authUC.Login(r.Context(), bookclub.LoginParams{
		Email:    payload.Email,
		Password: payload.Password,
	})

	if err != nil {
		renderError(r, w, err, "")

		return
	}

	resp := schema.LoginToken{
		Access:  access,
		Refresh: refresh,
	}

	renderJSON(r, w, http.StatusOK, resp)
}

func (b *bookClubAPI) PostAuthRefresh(w http.ResponseWriter, r *http.Request) {
	var payload schema.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		renderError(r, w, bookclub.ErrInvalidCredentials, bookclub.ErrInvalidBody.Error())

		return
	}

	access, refresh, err := b.authUC.Refresh(r.Context(), payload.Refresh)
	if err != nil {
		renderError(r, w, err, "")

		return
	}

	resp := schema.LoginToken{
		Access:  access,
		Refresh: refresh,
	}

	renderJSON(r, w, http.StatusOK, resp)
}
