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
		resp := schema.ErrorResponse{
			Code: schema.N401,
			Msg:  "invalid request body",
		}

		renderJSON(w, int(resp.Code), resp)

		return
	}

	access, refresh, err := b.authUC.Login(r.Context(), bookclub.LoginParams{
		Email:    payload.Email,
		Password: payload.Password,
	})

	if err != nil {
		var resp schema.ErrorResponse
		switch err {
		case bookclub.ErrInvalidCredentials:
			resp.Code = schema.N401
			resp.Msg = err.Error()
		default:
			resp.Code = schema.N500
			resp.Msg = bookclub.ErrInternalError.Error()
		}

		renderJSON(w, int(resp.Code), resp)

		return
	}

	resp := schema.LoginToken{
		Access:  access,
		Refresh: refresh,
	}

	renderJSON(w, http.StatusOK, resp)
}

func (b *bookClubAPI) PostAuthRefresh(w http.ResponseWriter, r *http.Request) {
	var payload schema.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		resp := schema.ErrorResponse{
			Code: schema.N401,
			Msg:  "invalid request body",
		}

		renderJSON(w, int(resp.Code), resp)

		return
	}

	access, refresh, err := b.authUC.Refresh(r.Context(), payload.Refresh)
	if err != nil {
		var resp schema.ErrorResponse
		switch err {
		case bookclub.ErrInvalidCredentials:
			resp.Code = schema.N401
			resp.Msg = err.Error()
		default:
			resp.Code = schema.N500
			resp.Msg = bookclub.ErrInternalError.Error()
		}

		renderJSON(w, int(resp.Code), resp)

		return
	}

	resp := schema.LoginToken{
		Access:  access,
		Refresh: refresh,
	}

	renderJSON(w, http.StatusOK, resp)
}
