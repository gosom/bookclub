package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/gosom/bookclub"
	"github.com/gosom/bookclub/api"
	"github.com/gosom/bookclub/api/schema"
	"github.com/gosom/bookclub/mocks"
)

func Test_PostUsers(t *testing.T) {
	mctrl := gomock.NewController(t)

	userUc := mocks.NewMockUserUseCases(mctrl)

	bc := api.NewBooklubAPI(userUc, nil)

	bodytpl := `{"email": "%s", "password": "%s"}`

	t.Run("with empty body", func(t *testing.T) {
		body := bytes.NewBufferString(``)

		req, err := http.NewRequest(http.MethodPost, "/users", body)
		require.NoError(t, err)

		rec := httptest.NewRecorder()

		bc.PostUsers(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.Equal(t, rec.Header().Get("Content-Type"), "application/json")

		var resp schema.ErrorResponse
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))

		require.Equal(t, schema.ErrorResponse{
			Code: schema.N400,
			Msg:  "invalid request body",
		}, resp)
	})

	t.Run("with invalid body", func(t *testing.T) {
		body := bytes.NewBufferString(`{not valid json}`)

		req, err := http.NewRequest(http.MethodPost, "/users", body)
		require.NoError(t, err)

		rec := httptest.NewRecorder()

		bc.PostUsers(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.Equal(t, rec.Header().Get("Content-Type"), "application/json")

		var resp schema.ErrorResponse
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))

		require.Equal(t, schema.ErrorResponse{
			Code: schema.N400,
			Msg:  "invalid request body",
		}, resp)
	})

	t.Run("with invalid password", func(t *testing.T) {
		userUc.EXPECT().
			Register(gomock.Any(), bookclub.RegisterParams{
				Email:    "john@doe.com",
				Password: "invalid",
			}).Return(bookclub.User{}, bookclub.ErrInvalidPassword)

		body := bytes.NewBufferString(fmt.Sprintf(bodytpl, "john@doe.com", "invalid"))

		req, err := http.NewRequest(http.MethodPost, "/users", body)
		require.NoError(t, err)

		rec := httptest.NewRecorder()

		bc.PostUsers(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.Equal(t, rec.Header().Get("Content-Type"), "application/json")

		var resp schema.ErrorResponse
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))

		require.Equal(t, schema.ErrorResponse{
			Code: schema.N400,
			Msg:  bookclub.ErrInvalidPassword.Error(),
		}, resp)
	})

	t.Run("with invalid email", func(t *testing.T) {
		userUc.EXPECT().
			Register(gomock.Any(), bookclub.RegisterParams{
				Email:    "john",
				Password: "123abc!A#4%",
			}).Return(bookclub.User{}, bookclub.ErrInvalidEmail)

		body := bytes.NewBufferString(fmt.Sprintf(bodytpl, "john", "123abc!A#4%"))

		req, err := http.NewRequest(http.MethodPost, "/users", body)
		require.NoError(t, err)

		rec := httptest.NewRecorder()

		bc.PostUsers(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.Equal(t, rec.Header().Get("Content-Type"), "application/json")

		var resp schema.ErrorResponse
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))

		require.Equal(t, schema.ErrorResponse{
			Code: schema.N400,
			Msg:  bookclub.ErrInvalidEmail.Error(),
		}, resp)
	})

	t.Run("with internal error", func(t *testing.T) {
		userUc.EXPECT().
			Register(gomock.Any(), gomock.Any()).
			Return(bookclub.User{}, bookclub.ErrInternalError)

		body := bytes.NewBufferString(fmt.Sprintf(bodytpl, "john", "123abc!A#4%"))

		req, err := http.NewRequest(http.MethodPost, "/users", body)
		require.NoError(t, err)

		rec := httptest.NewRecorder()

		bc.PostUsers(rec, req)

		require.Equal(t, http.StatusInternalServerError, rec.Code)
		require.Equal(t, rec.Header().Get("Content-Type"), "application/json")

		var resp schema.ErrorResponse
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))

		require.Equal(t, schema.ErrorResponse{
			Code: schema.N500,
			Msg:  bookclub.ErrInternalError.Error(),
		}, resp)
	})

	t.Run("with success", func(t *testing.T) {
		expectedUUID := uuid.New()
		expectedEmail := bookclub.Email("john@doe.com")
		expectedCreatedAt := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		expectedUpdatedAt := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

		userUc.EXPECT().
			Register(gomock.Any(), gomock.Any()).
			Return(bookclub.User{
				ID:        expectedUUID,
				Email:     expectedEmail,
				CreatedAt: expectedCreatedAt,
				UpdatedAt: expectedUpdatedAt,
			}, nil)

		body := bytes.NewBufferString(fmt.Sprintf(bodytpl, expectedEmail, "123abc!A#4%"))

		req, err := http.NewRequest(http.MethodPost, "/users", body)
		require.NoError(t, err)

		rec := httptest.NewRecorder()

		bc.PostUsers(rec, req)

		require.Equal(t, http.StatusCreated, rec.Code)
		require.Equal(t, rec.Header().Get("Content-Type"), "application/json")

		var resp schema.CreateUserResponse
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))

		require.Equal(t, schema.CreateUserResponse{
			Id:    expectedUUID.String(),
			Email: string(expectedEmail),
		}, resp)
	})

	t.Run("when email exists", func(t *testing.T) {
		userUc.EXPECT().
			Register(gomock.Any(), gomock.Any()).
			Return(bookclub.User{}, bookclub.ErrAlreadyExists)

		body := bytes.NewBufferString(fmt.Sprintf(bodytpl, "john@doe.com", "123abc!A#4%"))

		req, err := http.NewRequest(http.MethodPost, "/users", body)
		require.NoError(t, err)

		rec := httptest.NewRecorder()

		bc.PostUsers(rec, req)

		require.Equal(t, http.StatusConflict, rec.Code)

		var resp schema.ErrorResponse
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))

		require.Equal(t, schema.ErrorResponse{
			Code: schema.N409,
			Msg:  bookclub.ErrAlreadyExists.Error(),
		}, resp)
	})
}
