package dto

import (
	"time"

	"github.com/google/uuid"
)

type NewsItem struct {
	ID             uuid.UUID `json:"id"`
	ChannelID      uuid.UUID `json:"channel_id"`
	ChannelTitle   string    `json:"channel_title,omitempty"`
	ChannelSource  string    `json:"source_name,omitempty"`
	Title          string    `json:"title"`
	Link           string    `json:"link"`
	Description    *string   `json:"description,omitempty"`
	PubDate        time.Time `json:"pub_date"`
	FetchedAt      time.Time `json:"fetched_at"`
	RelevanceScore float64   `json:"relevance_score"`
	ImageURL       *string   `json:"image_url,omitempty"`
	// EnclosureURL   *string    `json:"enclosure_url,omitempty"`
}

// NewsListResponse — структура ответа списка
type NewsListResponse struct {
	Total int64      `json:"total"`
	Items []NewsItem `json:"items"`
	//Page       int        `json:"page"`
	//Limit      int        `json:"limit"`
	//TotalPages int        `json:"total_pages,omitempty"`
}
