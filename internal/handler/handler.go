package handler

import (
	"net/http"
	"time"
	"txsystem/internal/models"
)

type App interface {
	ChangeBalance(tx *models.Transaction) error
	CreateUser(user *models.User) error
}

type Handler struct {
	app App
}

func New(app App) *Handler {
	return &Handler{app: app}
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	tx := &models.Transaction{}
	tx.CreateAt = time.Now()

	if err := GetBody(r.Body, tx); err != nil {
		RespErr(w, err)
		return
	}

	if err := h.app.ChangeBalance(tx); err != nil {
		RespErr(w, err)
		return
	}

	RespOK(w, "")
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}

	if err := GetBody(r.Body, user); err != nil {
		RespErr(w, err)
		return
	}

	RespOK(w, "")
}
