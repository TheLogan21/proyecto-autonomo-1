package handlers

import (
	"net/http"

	"ebook-system/internal/middleware"
	"ebook-system/internal/services"
)

type AuthHandler struct {
	userService *services.UserService
}

func NewAuthHandler(us *services.UserService) *AuthHandler {
	return &AuthHandler{userService: us}
}

func (h *AuthHandler) LoginView(w http.ResponseWriter, r *http.Request) {
	render(w, r, "login", nil)
}

func (h *AuthHandler) LoginPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	user, err := h.userService.Login(r.Context(), r.FormValue("email"), r.FormValue("password"))
	if err != nil {
		render(w, r, "login", map[string]string{"Error": err.Error()})
		return
	}

	token := middleware.CreateSession(user.ID)
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
	})
	http.Redirect(w, r, "/catalog", http.StatusSeeOther)
}

func (h *AuthHandler) RegisterView(w http.ResponseWriter, r *http.Request) {
	render(w, r, "register", nil)
}

func (h *AuthHandler) RegisterPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	err := h.userService.Register(r.Context(), r.FormValue("username"), r.FormValue("email"), r.FormValue("password"))
	if err != nil {
		render(w, r, "register", map[string]string{"Error": err.Error()})
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("session"); err == nil {
		middleware.DestroySession(cookie.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
