package v1

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
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

func RespOK(c *fiber.Ctx, data any) {
	r := ResponceOK{
		Status:  http.StatusOK,
		Success: "true",
		Data:    data,
	}

	c.JSON(r)
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
