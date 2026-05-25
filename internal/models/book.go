package models

import "time"

type Book struct {
	ID            int       `json:"id"`
	Title         string    `json:"title"`
	Author        string    `json:"author"`
	PublishedYear int       `json:"publishedYear"`
	ISBN          string    `json:"isbn"`
	CreatedAt     time.Time `json:"createdAt"`
}
