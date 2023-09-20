package api

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func renderJSON(w http.ResponseWriter, code int, body any) {
	w.Header().Set("Content-Type", "application/json")

	buff := bytes.NewBuffer(nil)

	if err := json.NewEncoder(buff).Encode(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"code": 500, "msg": "internal server error"}`))

		return
	}

	w.WriteHeader(code)
	_, _ = w.Write(buff.Bytes())
}
