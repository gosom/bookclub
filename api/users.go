package api

import (
	"encoding/json"
	"net/http"

	"github.com/gosom/bookclub"
	"github.com/gosom/bookclub/api/schema"
)

func (b *bookClubAPI) PostUsers(w http.ResponseWriter, r *http.Request) {
	var payload schema.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		renderError(r, w, bookclub.ErrInvalidBody, "")

		return
	}

	user, err := b.userUC.Register(r.Context(), bookclub.RegisterParams{
		Email:    payload.Email,
		Password: payload.Password,
	})
	if err != nil {
		renderError(r, w, err, "")

		return
	}

	ans := schema.CreateUserResponse{
		Email: string(user.Email),
		Id:    user.ID.String(),
	}

	renderJSON(r, w, http.StatusCreated, ans)
}
