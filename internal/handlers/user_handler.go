package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"ebook-system/internal/middleware"
	"ebook-system/internal/services"
)

type UserHandler struct {
	userService     services.UserService
	purchaseService services.PurchaseService
}

func NewUserHandler(us services.UserService, ps services.PurchaseService) *UserHandler {
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

func (h *UserHandler) ProfileJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := middleware.GetUserID(r.Context())

	user, err := h.userService.GetUserByID(r.Context(), userID)
	if err != nil || user == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Usuario no encontrado"})
		return
	}

	purchases, err := h.purchaseService.GetUserPurchases(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al obtener compras"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user":      user,
		"purchases": purchases,
	})
}

func (h *UserHandler) BalanceJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Método no permitido"})
		return
	}

	userID := middleware.GetUserID(r.Context())

	var req struct {
		Amount float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Cuerpo JSON inválido o malformado"})
		return
	}

	if req.Amount <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "La cantidad debe ser mayor a cero"})
		return
	}

	err := h.userService.AddBalance(r.Context(), userID, req.Amount)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Saldo actualizado exitosamente"})
}

