package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

type ResponceOK struct {
	Status  int         `json:"status,omitempty"`
	Success string      `json:"success,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func RespErr(w http.ResponseWriter, err error) {
	enc := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var res struct {
		Err string `json:"error"`
	}
	res.Err = err.Error()

	if err = enc.Encode(res); err != nil {
		http.Error(w, "failed encode responce", http.StatusInternalServerError)
	}
}

func RespOK(w http.ResponseWriter, data interface{}) {
	r := ResponceOK{
		Status:  http.StatusOK,
		Success: "true",
		Data:    data,
	}

	enc := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := enc.Encode(r); err != nil {
		http.Error(w, "failed encode responce", http.StatusInternalServerError)
	}
}

func GetBody(in io.ReadCloser, out interface{}) error {
	if in == nil {
		return errors.New("nil request body")
	}

	defer in.Close()
	buf, err := io.ReadAll(in)
	if err != nil {
		return errors.Wrap(err, "failed read body")
	}

	err = json.Unmarshal(buf, out)
	if err != nil {
		return errors.Wrap(err, "failed unmarshal body")
	}

	return nil
}
