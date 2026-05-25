package handlers

import (
	"net/http"
	"net/url"
	"strconv"

	"ebook-system/internal/middleware"
	"ebook-system/internal/services"
)

type UserHandler struct {
	userService     *services.UserService
	purchaseService *services.PurchaseService
}

func NewUserHandler(us *services.UserService, ps *services.PurchaseService) *UserHandler {
	return &UserHandler{userService: us, purchaseService: ps}
}

func (h *UserHandler) Profile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	user, err := h.userService.GetUserByID(r.Context(), userID)
	if err != nil || user == nil {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	purchases, err := h.purchaseService.GetUserPurchases(r.Context(), userID)
	if err != nil {
		http.Error(w, "Error al obtener compras", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"User":      user,
		"Purchases": purchases,
	}

	render(w, r, "profile", data)
}

func (h *UserHandler) AddBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/profile?error="+url.QueryEscape("Método no permitido"), http.StatusSeeOther)
		return
	}

	userID := middleware.GetUserID(r.Context())

	r.ParseForm()
	amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
	if err != nil || amount <= 0 {
		http.Redirect(w, r, "/profile?error="+url.QueryEscape("Cantidad inválida"), http.StatusSeeOther)
		return
	}

	err = h.userService.AddBalance(r.Context(), userID, amount)
	if err != nil {
		http.Redirect(w, r, "/profile?error="+url.QueryEscape(err.Error()), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
