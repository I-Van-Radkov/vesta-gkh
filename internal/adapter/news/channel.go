package news

import (
	"context"
	"fmt"
	"time"

	"github.com/I-Van-Radkov/vesta-gkh/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *NewsRepo) GetChannelById(id uuid.UUID) (*models.Channel, error) {
	var ch models.Channel
	err := r.db.QueryRow(context.Background(),
		`SELECT id, title, url, source_name, description, language,
		        last_build_date, generator, image_url, is_active, priority,
		        created_at, updated_at
		 FROM channels
		 WHERE id = $1`,
		id,
	).Scan(
		&ch.ID, &ch.Title, &ch.URL, &ch.SourceName, &ch.Description,
		&ch.Language, &ch.LastBuildDate, &ch.Generator, &ch.ImageURL,
		&ch.IsActive, &ch.Priority, &ch.CreatedAt, &ch.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get channel by id %s: %w", id, err)
	}
	return &ch, nil
}

func (r *NewsRepo) GetActiveChannels() ([]*models.Channel, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, title, url, source_name, description, language,
		        last_build_date, generator, image_url, is_active, priority,
		        created_at, updated_at
		 FROM channels
		 WHERE is_active = true
		 ORDER BY priority ASC, created_at ASC`)
	if err != nil {
		return nil, fmt.Errorf("query active channels: %w", err)
	}
	defer rows.Close()

	var channels []*models.Channel
	for rows.Next() {
		var ch models.Channel
		err := rows.Scan(
			&ch.ID, &ch.Title, &ch.URL, &ch.SourceName, &ch.Description,
			&ch.Language, &ch.LastBuildDate, &ch.Generator, &ch.ImageURL,
			&ch.IsActive, &ch.Priority, &ch.CreatedAt, &ch.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan channel: %w", err)
		}
		channels = append(channels, &ch)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return channels, nil
}

func (r *NewsRepo) GetURLsChannels(id uuid.UUID) ([]string, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT url FROM channels WHERE is_active = true`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []string
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}
	return urls, rows.Err()
}

func (r *NewsRepo) GetLatestPubDateByChannel(ctx context.Context, channelID uuid.UUID) (time.Time, error) {
	var t time.Time
	err := r.db.QueryRow(ctx,
		`SELECT COALESCE(MAX(pub_date), '0001-01-01'::timestamptz)
		 FROM news
		 WHERE channel_id = $1`,
		channelID,
	).Scan(&t)

	if err != nil {
		return time.Time{}, fmt.Errorf("get latest pub_date for channel %s: %w", channelID, err)
	}
	return t, nil
}
