package api

import (
	"github.com/gosom/bookclub"
	"github.com/gosom/bookclub/api/schema"
)

var _ schema.ServerInterface = (*bookClubAPI)(nil)

type bookClubAPI struct {
	userUC bookclub.UserUseCases
	authUC bookclub.AuthUseCases
}

func NewBooklubAPI(
	userUC bookclub.UserUseCases,
	authUC bookclub.AuthUseCases,
) schema.ServerInterface {

	return &bookClubAPI{
		userUC: userUC,
		authUC: authUC,
	}
}
