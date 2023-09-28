package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/gosom/bookclub"
	"github.com/gosom/bookclub/api/schema"
	"github.com/gosom/bookclub/observ"
)

func renderJSON(r *http.Request, w http.ResponseWriter, code int, body any) {
	w.Header().Set("Content-Type", "application/json")

	buff := bytes.NewBuffer(nil)

	if err := json.NewEncoder(buff).Encode(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"code": 500, "msg": "internal server error"}`))

		setReqCtxValue(r, observ.ErrorCtxKey, err)

		return
	}

	w.WriteHeader(code)
	_, _ = w.Write(buff.Bytes())

	we, ok := body.(wrappedErrorResponse)
	if ok {
		setReqCtxValue(r, observ.ErrorCtxKey, we.err)
	}
}

func renderError(r *http.Request, w http.ResponseWriter, err error, msg string) {
	if err == nil {
		err = bookclub.ErrInternalError
	}

	var code schema.ErrorResponseCode

	switch err {
	case bookclub.ErrInvalidCredentials:
		code = schema.N401
	case bookclub.ErrInvalidEmail, bookclub.ErrInvalidPassword:
		code = schema.N400
	case bookclub.ErrAlreadyExists:
		code = schema.N409
	case bookclub.ErrInvalidBody:
		code = schema.N400
	default:
		code = schema.N500
		if msg == "" {
			msg = bookclub.ErrInternalError.Error()
		}
	}

	resp := newErrorResponse(err, schema.ErrorResponseCode(code), msg)

	renderJSON(r, w, int(code), resp)
}

func setReqCtxValue(r *http.Request, key, val any) {
	ctx := contextWithValue(r.Context(), key, val)

	*r = *r.WithContext(ctx)
}

func contextWithValue(ctx context.Context, key, val any) context.Context {
	return context.WithValue(ctx, key, val)
}

type wrappedErrorResponse struct {
	err error
	schema.ErrorResponse
}

func newErrorResponse(err error, code schema.ErrorResponseCode, msg string) wrappedErrorResponse {
	if msg == "" {
		msg = err.Error()
	}
	return wrappedErrorResponse{
		err: err,
		ErrorResponse: schema.ErrorResponse{
			Code: code,
			Msg:  msg,
		},
	}
}
