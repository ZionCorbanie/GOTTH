package handlers

import (
	b64 "encoding/base64"
	"fmt"
	"net/http"
	"goth/internal/hash"
	"goth/internal/store"
	"goth/internal/templates"
	"time"
)

type LoginHandler struct {
	userStore         store.UserStore
	sessionStore      store.SessionStore
	passwordhash      hash.PasswordHash
	sessionCookieName string
}

type LoginHandlerParams struct {
	UserStore         store.UserStore
	SessionStore      store.SessionStore
	PasswordHash      hash.PasswordHash
	SessionCookieName string
}

func NewLoginHandler(params LoginHandlerParams) *LoginHandler {
	return &LoginHandler{
		userStore:         params.UserStore,
		sessionStore:      params.SessionStore,
		passwordhash:      params.PasswordHash,
		sessionCookieName: params.SessionCookieName,
	}
}

func (h *LoginHandler) Show(w http.ResponseWriter, r *http.Request) {
    target := r.URL.Query().Get("redirect")

    if target == "" {
        target = "/"
    }

	c := templates.Login("Login", target)
	err := templates.Layout(c, "Robokar - Login").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func (h *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {

	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := h.userStore.GetUser(username)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		c := templates.LoginError()
		c.Render(r.Context(), w)
		return
	}

	passwordIsValid, err := h.passwordhash.ComparePasswordAndHash(password, user.Password)

	if err != nil || !passwordIsValid {
		w.WriteHeader(http.StatusUnauthorized)
		c := templates.LoginError()
		c.Render(r.Context(), w)
		return
	}

	session, err := h.sessionStore.CreateSession(&store.Session{
		UserID: user.ID,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userID := user.ID
	sessionID := session.SessionID

	cookieValue := b64.StdEncoding.EncodeToString(fmt.Appendf(nil, "%s:%d", sessionID, userID))

	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{
		Name:     h.sessionCookieName,
		Value:    cookieValue,
		Expires:  expiration,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &cookie)

    target := r.URL.Query().Get("redirect")

	w.Header().Set("HX-Redirect", target)
	w.WriteHeader(http.StatusOK)
}

func (h *LoginHandler) Logout(w http.ResponseWriter, r *http.Request) {

	http.SetCookie(w, &http.Cookie{
		Name:    h.sessionCookieName,
		MaxAge:  -1,
		Expires: time.Now().Add(-100 * time.Hour),
		Path:    "/",
	})

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}
