package handlers

import (
	"goth/internal/store"
	"goth/internal/templates"
	"net/http"
)


type RegisterHandler struct {
	userStore store.UserStore
}

type RegisterHandlerParams struct {
	UserStore store.UserStore
}

func NewRegisterHandler(params RegisterHandlerParams) *RegisterHandler {
	return &RegisterHandler{
		userStore: params.UserStore,
	}
}

func (h *RegisterHandler) Show(w http.ResponseWriter, r *http.Request) {
	c := templates.RegisterPage()
	err := templates.Layout(c, "My website").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}

}

func (h *RegisterHandler) Register(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	err := h.userStore.CreateUser(email, password)

	if err != nil {

		w.WriteHeader(http.StatusBadRequest)
		c := templates.RegisterError()
		c.Render(r.Context(), w)
		return
	}

	c := templates.RegisterSuccess()
	err = c.Render(r.Context(), w)

	if err != nil {
		http.Error(w, "error rendering template", http.StatusInternalServerError)
		return
	}

}
