package api

import (
	"encoding/json"
	"net/http"

	"github.com/gosom/bookclub/api/schema"
)

var _ schema.ServerInterface = (*bookClubAPI)(nil)

type bookClubAPI struct {
}

func NewBooklubAPI() schema.ServerInterface {
	return &bookClubAPI{}
}

func (b *bookClubAPI) PostUsers(w http.ResponseWriter, r *http.Request) {
	var payload schema.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		resp := schema.ErrorResponse{
			Code: schema.N400,
			Msg:  "invalid request body",
		}

		renderJSON(w, int(resp.Code), resp)

		return
	}

	// TODO here we need to call the service to create the user
	// and return him. We skip for now. Email and passwords needs
	// to be also validated.

	ans := schema.CreateUserResponse{
		Email: payload.Email,
		Id:    "c69625f8-57ab-11ee-8d34-38f3ab390182",
	}

	renderJSON(w, http.StatusCreated, ans)
}
