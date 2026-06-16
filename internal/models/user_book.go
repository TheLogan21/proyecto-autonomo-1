package models

import "time"

type UserBook struct {
	UserID        int       `json:"userId"`
	BookID        int       `json:"bookId"`
	DownloadCount int       `json:"downloadCount"`
	PurchasedAt   time.Time `json:"purchasedAt"`
	// Detalle del libro relacionado
	Book *Book `json:"book,omitempty"`
}
