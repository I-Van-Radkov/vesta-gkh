package models

import (
	"encoding/json"
	"regexp"
	"time"

	"github.com/google/uuid"
)

type Channel struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	Title         string     `json:"title" db:"title"`
	URL           string     `json:"url" db:"url"`
	SourceName    string     `json:"source_name,omitempty" db:"source_name"`
	Description   *string    `json:"description" db:"description"`
	Language      string     `json:"language" db:"language"`
	LastBuildDate *time.Time `json:"last_build_date,omitempty" db:"last_build_date"`
	Generator     *string    `json:"generator,omitempty" db:"generator"`
	ImageURL      *string    `json:"image_url,omitempty" db:"image_url"`
	IsActive      bool       `json:"is_active" db:"is_active"`
	Priority      int        `json:"priority" db:"priority"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

type News struct {
	ID              uuid.UUID       `json:"id" db:"id"`
	ChannelID       uuid.UUID       `json:"channel_id" db:"channel_id"`
	GUID            *string         `json:"guid,omitempty" db:"guid"`
	Title           string          `json:"title" db:"title"`
	Link            string          `json:"link" db:"link"`
	Description     *string         `json:"description,omitempty" db:"description"`
	FullText        *string         `json:"full_text,omitempty" db:"full_text"`
	PubDate         time.Time       `json:"pub_date" db:"pub_date"`
	FetchedAt       time.Time       `json:"fetched_at" db:"fetched_at"`
	EnclosureURL    *string         `json:"enclosure_url,omitempty" db:"enclosure_url"`
	EnclosureType   *string         `json:"enclosure_type,omitempty" db:"enclosure_type"`
	EnclosureLength *int64          `json:"enclosure_length,omitempty" db:"enclosure_length"`
	MediaGroup      json.RawMessage `json:"media_group,omitempty" db:"media_group"`
	RelevanceScore  float64         `json:"relevance_score" db:"relevance_score"`
	CreatedAt       time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at" db:"updated_at"`
}

type GKHKeyword struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Keyword   string    `json:"keyword" db:"keyword"`
	Weight    float64   `json:"weight" db:"weight"`
	Category  *string   `json:"category" db:"category"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type GKHRegexRule struct {
	ID         uuid.UUID      `json:"id" db:"id"`
	PatternStr string         `json:"pattern" db:"pattern"`
	Compiled   *regexp.Regexp `json:"-" db:"-"`
	BonusScore float64        `json:"bonus_score" db:"bonus_score"`
	IsActive   bool           `json:"is_active" db:"is_active"`
	CreatedAt  time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at" db:"updated_at"`
}

type GKHNegativeRegex struct {
	ID          int64          `json:"id" db:"id"`
	PatternStr  string         `json:"pattern" db:"pattern"`
	Compiled    *regexp.Regexp `json:"-" db:"-"`
	Penalty     float64        `json:"penalty" db:"penalty"`
	Description *string        `json:"description" db:"description"`
	IsActive    bool           `json:"is_active" db:"is_active"`
	CreatedAt   time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at" db:"updated_at"`
}

type GHKCategoryWithKeywords struct {
	ID              uuid.UUID `json:"id" db:"id"`
	Name            string    `json:"name" db:"name"`
	BonusPerHit     float64   `json:"bonus_per_hit" db:"bonus_per_hit"`
	MinHitsForBonus int       `json:"min_hits_for_bonus" db:"min_hits_for_bonus"`
	BigBonus        float64   `json:"big_bonus" db:"big_bonus"`
	Description     *string   `json:"description" db:"description"`
	Keywords        []string  `json:"keywords"`
}
