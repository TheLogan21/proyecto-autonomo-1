package models

import "time"

type UserBook struct {
	UserID        int       `json:"userId"`
	BookID        int       `json:"bookId"`
	DownloadCount int       `json:"downloadCount"`
	PurchasedAt   time.Time `json:"purchasedAt"`
	// Relación para mostrar detalles del libro en el perfil
	Book *Book `json:"book,omitempty"`
}
