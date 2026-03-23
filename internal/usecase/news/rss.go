package news

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/I-Van-Radkov/vesta-gkh/internal/models"
	"github.com/google/uuid"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title         string    `xml:"title"`
	Link          string    `xml:"link"`
	Description   string    `xml:"description"`
	Language      string    `xml:"language"`
	PubDate       string    `xml:"pubDate"`
	LastBuildDate string    `xml:"lastBuildDate"`
	Generator     string    `xml:"generator"`
	Items         []RSSItem `xml:"item"`
	Image         *struct {
		URL   string `xml:"url"`
		Title string `xml:"title"`
		Link  string `xml:"link"`
	} `xml:"image"`
}

type RSSItem struct {
	Title       string      `xml:"title"`
	Link        string      `xml:"link"`
	GUID        string      `xml:"guid"`
	Description string      `xml:"description"`
	PubDate     string      `xml:"pubDate"`
	FullText    string      `xml:"yandex:full-text"`
	Encoded     string      `xml:"encoded"`
	Content     string      `xml:"content"`
	Enclosure   *Enclosure  `xml:"enclosure"`
	Media       *MediaGroup `xml:"media:group"`
}

type Enclosure struct {
	URL    string `xml:"url,attr"`
	Length int64  `xml:"length,attr"`
	Type   string `xml:"type,attr"`
}

type MediaGroup struct {
	Content []MediaContent `xml:"media:content"`
}

type MediaContent struct {
	URL    string `xml:"url,attr"`
	Type   string `xml:"type,attr"`
	Length int64  `xml:"length,attr"`
}

func (p *NewsParser) ParseActiveChannels(ctx context.Context) ([]*models.News, error) {
	channels, err := p.newsRepo.GetActiveChannels()
	if err != nil {
		return nil, fmt.Errorf("get active channels: %w", err)
	}

	var allNews []*models.News
	for _, ch := range channels {
		items, err := p.ParseChannel(ctx, ch)
		if err != nil {
			log.Printf("парсинг канала %s (%s) упал: %v", ch.Title, ch.URL, err)
			continue
		}
		allNews = append(allNews, items...)
	}
	return allNews, nil
}

func (p *NewsParser) ParseChannel(ctx context.Context, ch *models.Channel) ([]*models.News, error) {
	latest, _ := p.newsRepo.GetLatestPubDateByChannel(ctx, ch.ID)

	resp, err := p.client.Get(ch.URL)
	if err != nil {
		return nil, fmt.Errorf("не удалось загрузить RSS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("получен статус: %s", resp.Status)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 10*1024*1024))
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать данные: %w", err)
	}

	var rss RSS
	if err := xml.Unmarshal(body, &rss); err != nil {
		return nil, fmt.Errorf("не удалось распарсить XML: %w", err)
	}

	var newsList []*models.News
	now := time.Now()

	for _, item := range rss.Channel.Items {
		pubDate := p.parsePubDate(item.PubDate, now)

		// Если новость старше latest - 60 минут, пропускаем
		if !pubDate.After(latest.Add(-60 * time.Minute)) {
			continue
		}

		fullText := p.extractFullText(item)

		title := clean(item.Title)
		desc := clean(item.Description)
		content := clean(fullText)

		ok, score, debug := p.classify(title, desc, content)
		if score > p.cfg.Threshold {
			log.Printf("[NEWS] score=%.1f | %s | %s", score, title, debug)
		}
		if !ok {
			continue
		}

		guid := item.GUID
		if guid == "" && item.Link != "" {
			guid = item.Link
		}

		news := &models.News{
			ID:             uuid.New(),
			ChannelID:      ch.ID,
			GUID:           ptrIfNotEmpty(guid),
			Title:          title,
			Link:           item.Link,
			Description:    ptrIfNotEmpty(desc),
			FullText:       ptrIfNotEmpty(content),
			PubDate:        pubDate,
			FetchedAt:      now,
			RelevanceScore: score,
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		newsList = append(newsList, news)
	}

	return newsList, nil
}
