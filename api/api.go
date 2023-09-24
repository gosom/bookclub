package api

import (
	"encoding/json"
	"net/http"

	"github.com/gosom/bookclub"
	"github.com/gosom/bookclub/api/schema"
)

var _ schema.ServerInterface = (*bookClubAPI)(nil)

type bookClubAPI struct {
	userUC bookclub.UserUseCases
}

func NewBooklubAPI(userUC bookclub.UserUseCases) schema.ServerInterface {
	return &bookClubAPI{
		userUC: userUC,
	}
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

	user, err := b.userUC.Register(r.Context(), bookclub.RegisterParams{
		Email:    payload.Email,
		Password: payload.Password,
	})
	if err != nil {
		resp := schema.ErrorResponse{}

		switch err {
		case bookclub.ErrInvalidEmail:
			resp.Code = schema.N400
			resp.Msg = err.Error()
		case bookclub.ErrInvalidPassword:
			resp.Code = schema.N400
			resp.Msg = err.Error()
		case bookclub.ErrInternalError:
			resp.Code = schema.N500
			resp.Msg = err.Error()
		case bookclub.ErrAlreadyExists:
			resp.Code = schema.N409
			resp.Msg = err.Error()
		default:
			resp.Code = schema.N500
			resp.Msg = "internal server error"
		}

		renderJSON(w, int(resp.Code), resp)

		return
	}

	ans := schema.CreateUserResponse{
		Email: string(user.Email),
		Id:    user.ID.String(),
	}

	renderJSON(w, http.StatusCreated, ans)
}
