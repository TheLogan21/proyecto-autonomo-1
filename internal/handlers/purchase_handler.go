package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"ebook-system/internal/middleware"
	"ebook-system/internal/services"
	"github.com/go-chi/chi/v5"
)

type PurchaseHandler struct {
	purchaseService services.PurchaseService
}

func NewPurchaseHandler(ps services.PurchaseService) *PurchaseHandler {
	return &PurchaseHandler{purchaseService: ps}
}

func (h *PurchaseHandler) BuyBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/catalog?error="+url.QueryEscape("Método no permitido"), http.StatusSeeOther)
		return
	}

	userID := middleware.GetUserID(r.Context())
	bookID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Redirect(w, r, "/catalog?error="+url.QueryEscape("ID inválido"), http.StatusSeeOther)
		return
	}

	err = h.purchaseService.BuyBook(r.Context(), userID, bookID)
	if err != nil {
		http.Redirect(w, r, "/catalog?error="+url.QueryEscape(err.Error()), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func (h *PurchaseHandler) DownloadBook(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	bookID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Redirect(w, r, "/profile?error="+url.QueryEscape("ID inválido"), http.StatusSeeOther)
		return
	}

	content, err := h.purchaseService.DownloadBook(r.Context(), userID, bookID)
	if err != nil {
		http.Redirect(w, r, "/profile?error="+url.QueryEscape(err.Error()), http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=book_download.txt")
	w.Header().Set("Content-Type", "text/plain")
	w.Write(content)
}

func (h *PurchaseHandler) BuyJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Método no permitido"})
		return
	}

	userID := middleware.GetUserID(r.Context())
	bookID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "ID de libro inválido"})
		return
	}

	err = h.purchaseService.BuyBook(r.Context(), userID, bookID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Libro comprado exitosamente"})
}

