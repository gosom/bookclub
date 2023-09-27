package api_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gosom/bookclub"
	"github.com/gosom/bookclub/api"
	"github.com/gosom/bookclub/api/schema"
	"github.com/gosom/bookclub/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_PostAuthLogin(t *testing.T) {
	mctrl := gomock.NewController(t)

	authUc := mocks.NewMockAuthUseCases(mctrl)

	bc := api.NewBooklubAPI(nil, authUc)

	bodyTpl := `{ "email": "%s", "password": "%s" }`

	t.Run("with valid credentials", func(t *testing.T) {
		expected := schema.LoginToken{
			Access:  "access",
			Refresh: "refresh",
		}

		authUc.EXPECT().Login(gomock.Any(), bookclub.LoginParams{
			Email:    "j@doe.com",
			Password: "123Ab1c$139afs%",
		}).Return(expected.Access, expected.Refresh, nil)

		body := bytes.NewBufferString(
			fmt.Sprintf(bodyTpl, "j@doe.com", "123Ab1c$139afs%"),
		)

		req, err := http.NewRequest(http.MethodPost, "/auth/login", body)
		require.NoError(t, err)

		rec := httptest.NewRecorder()

		bc.PostAuthLogin(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		require.Equal(t, rec.Header().Get("Content-Type"), "application/json")

		var resp schema.LoginToken
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))

		require.Equal(t, expected, resp)
	})

	t.Run("with invalid credentials", func(t *testing.T) {
		authUc.EXPECT().Login(gomock.Any(), bookclub.LoginParams{
			Email:    "j@doe.com",
			Password: "invalid",
		}).Return("", "", bookclub.ErrInvalidCredentials)

		body := bytes.NewBufferString(
			fmt.Sprintf(bodyTpl, "j@doe.com", "invalid"),
		)

		req, err := http.NewRequest(http.MethodPost, "/auth/login", body)
		require.NoError(t, err)

		rec := httptest.NewRecorder()

		bc.PostAuthLogin(rec, req)

		require.Equal(t, http.StatusUnauthorized, rec.Code)
		require.Equal(t, rec.Header().Get("Content-Type"), "application/json")

		expected := schema.ErrorResponse{
			Code: schema.N401,
			Msg:  bookclub.ErrInvalidCredentials.Error(),
		}

		var resp schema.ErrorResponse
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))

		require.Equal(t, expected, resp)
	})

	t.Run("with any error but ErrInvalidCredentials", func(t *testing.T) {
		authUc.EXPECT().Login(gomock.Any(), bookclub.LoginParams{
			Email:    "j@doe.com",
			Password: "invalid",
		}).Return("", "", errors.New("any error"))

		body := bytes.NewBufferString(
			fmt.Sprintf(bodyTpl, "j@doe.com", "invalid"),
		)

		req, err := http.NewRequest(http.MethodPost, "/auth/login", body)
		require.NoError(t, err)

		rec := httptest.NewRecorder()

		bc.PostAuthLogin(rec, req)

		require.Equal(t, http.StatusInternalServerError, rec.Code)
		require.Equal(t, rec.Header().Get("Content-Type"), "application/json")

		expected := schema.ErrorResponse{
			Code: schema.N500,
			Msg:  bookclub.ErrInternalError.Error(),
		}

		var resp schema.ErrorResponse
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))

		require.Equal(t, expected, resp)
	})
}

func Test_PostAuthRefresh(t *testing.T) {
	mctrl := gomock.NewController(t)

	authUc := mocks.NewMockAuthUseCases(mctrl)

	bc := api.NewBooklubAPI(nil, authUc)

	bodyTpl := `{ "refresh": "%s" }`

	t.Run("with invalid json", func(t *testing.T) {
		body := bytes.NewBufferString("not a json")
		req, err := http.NewRequest(http.MethodPost, "/auth/refresh", body)
		require.NoError(t, err)

		rec := httptest.NewRecorder()

		bc.PostAuthRefresh(rec, req)

		require.Equal(t, http.StatusUnauthorized, rec.Code)
		require.Equal(t, rec.Header().Get("Content-Type"), "application/json")

		expected := schema.ErrorResponse{
			Code: schema.N401,
			Msg:  "invalid request body",
		}

		var resp schema.ErrorResponse
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))

		require.Equal(t, expected, resp)
	})

	t.Run("with a valid refresh token", func(t *testing.T) {
		expected := schema.LoginToken{
			Access:  "access",
			Refresh: "refresh",
		}

		authUc.EXPECT().Refresh(gomock.Any(), "old refresh token").
			Return(expected.Access, expected.Refresh, nil)

		body := bytes.NewBufferString(
			fmt.Sprintf(bodyTpl, "old refresh token"),
		)

		req, err := http.NewRequest(http.MethodPost, "/auth/refresh", body)
		require.NoError(t, err)

		rec := httptest.NewRecorder()

		bc.PostAuthRefresh(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		require.Equal(t, rec.Header().Get("Content-Type"), "application/json")

		var resp schema.LoginToken
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))

		require.Equal(t, expected, resp)
	})

	t.Run("with an invalid refresh token", func(t *testing.T) {
		authUc.EXPECT().Refresh(gomock.Any(), "invalid refresh token").
			Return("", "", bookclub.ErrInvalidCredentials)

		body := bytes.NewBufferString(
			fmt.Sprintf(bodyTpl, "invalid refresh token"),
		)
		req, err := http.NewRequest(http.MethodPost, "/auth/refresh", body)
		require.NoError(t, err)

		rec := httptest.NewRecorder()

		bc.PostAuthRefresh(rec, req)

		require.Equal(t, http.StatusUnauthorized, rec.Code)
		require.Equal(t, rec.Header().Get("Content-Type"), "application/json")

		expected := schema.ErrorResponse{
			Code: schema.N401,
			Msg:  bookclub.ErrInvalidCredentials.Error(),
		}

		var resp schema.ErrorResponse
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))

		require.Equal(t, expected, resp)
	})
}
