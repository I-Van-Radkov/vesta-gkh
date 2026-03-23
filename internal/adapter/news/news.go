package news

import (
	"context"
	"fmt"
	"log"

	"github.com/I-Van-Radkov/vesta-gkh/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *NewsRepo) GetNewsById(id uuid.UUID) (*models.News, error) {
	var n models.News
	err := r.db.QueryRow(context.Background(),
		`SELECT id, channel_id, guid, title, link, description, full_text,
		        pub_date, fetched_at, enclosure_url, enclosure_type, enclosure_length,
		        media_group, relevance_score, created_at, updated_at
		 FROM news
		 WHERE id = $1`,
		id,
	).Scan(
		&n.ID, &n.ChannelID, &n.GUID, &n.Title, &n.Link, &n.Description, &n.FullText,
		&n.PubDate, &n.FetchedAt, &n.EnclosureURL, &n.EnclosureType, &n.EnclosureLength,
		&n.MediaGroup, &n.RelevanceScore, &n.CreatedAt, &n.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get news by id %s: %w", id, err)
	}
	return &n, nil
}

func (r *NewsRepo) GetNewsList() ([]*models.News, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, channel_id, guid, title, link, description, full_text,
		        pub_date, fetched_at, enclosure_url, enclosure_type, enclosure_length,
		        media_group, relevance_score, created_at, updated_at
		 FROM news
		 ORDER BY pub_date DESC
		 LIMIT 100`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*models.News
	for rows.Next() {
		var n models.News
		err := rows.Scan(
			&n.ID, &n.ChannelID, &n.GUID, &n.Title, &n.Link, &n.Description, &n.FullText,
			&n.PubDate, &n.FetchedAt, &n.EnclosureURL, &n.EnclosureType, &n.EnclosureLength,
			&n.MediaGroup, &n.RelevanceScore, &n.CreatedAt, &n.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, &n)
	}
	return list, rows.Err()
}

func (r *NewsRepo) AddNews(n *models.News) (uuid.UUID, error) {
	var insertedID uuid.UUID
	err := r.db.QueryRow(context.Background(),
		`INSERT INTO news (
			id, channel_id, guid, title, link, description, full_text,
			pub_date, fetched_at, enclosure_url, enclosure_type, enclosure_length,
			media_group, relevance_score, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		ON CONFLICT (link) DO UPDATE SET
			title         = EXCLUDED.title,
			description   = EXCLUDED.description,
			full_text     = EXCLUDED.full_text,
			pub_date      = EXCLUDED.pub_date,
			fetched_at    = EXCLUDED.fetched_at,
			relevance_score = EXCLUDED.relevance_score,
			updated_at    = EXCLUDED.updated_at
		RETURNING id`,
		n.ID, n.ChannelID, n.GUID, n.Title, n.Link, n.Description, n.FullText,
		n.PubDate, n.FetchedAt, n.EnclosureURL, n.EnclosureType, n.EnclosureLength,
		n.MediaGroup, n.RelevanceScore, n.CreatedAt, n.UpdatedAt,
	).Scan(&insertedID)

	if err != nil {
		return uuid.Nil, fmt.Errorf("insert/update news %s: %w", n.Link, err)
	}
	return insertedID, nil
}

func (r *NewsRepo) AddNewsList(items []*models.News) ([]uuid.UUID, error) {
	if len(items) == 0 {
		return nil, nil
	}

	ctx := context.Background()
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var ids []uuid.UUID

	for _, n := range items {
		id, err := r.AddNews(n) // можно было бы сделать батч, но для простоты по одному
		if err != nil {
			log.Printf("failed to add news %s: %v", n.Link, err)
			continue
		}
		ids = append(ids, id)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return ids, nil
}
