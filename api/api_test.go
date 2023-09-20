package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gosom/bookclub/api"
	"github.com/gosom/bookclub/api/schema"
)

func Test_PostUsers(t *testing.T) {
	bc := api.NewBooklubAPI()

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

	t.Run("with a valid json", func(t *testing.T) {
		body := bytes.NewBufferString(`{"email": "john@doe.com"}`)

		req, err := http.NewRequest(http.MethodPost, "/users", body)
		require.NoError(t, err)

		rec := httptest.NewRecorder()

		bc.PostUsers(rec, req)

		require.Equal(t, http.StatusCreated, rec.Code)
		require.Equal(t, rec.Header().Get("Content-Type"), "application/json")

		expectedResponse := schema.CreateUserResponse{
			Email: "john@doe.com",
			Id:    "c69625f8-57ab-11ee-8d34-38f3ab390182",
		}

		var resp schema.CreateUserResponse
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))

		require.Equal(t, expectedResponse, resp)
	})
}
