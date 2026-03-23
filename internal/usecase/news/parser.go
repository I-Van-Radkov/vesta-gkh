package news

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/I-Van-Radkov/vesta-gkh/internal/models"
)

type NewsParser struct {
	client   *http.Client
	newsRepo NewsAdapterProvider

	cfg ParserConfig

	mu         sync.Mutex
	keywords   []*models.GKHKeyword
	regexRules []*models.GKHRegexRule
	negRules   []*models.GKHNegativeRegex
	categories []*models.GHKCategoryWithKeywords
	lastCache  time.Time
}

func NewNewsParser(repo NewsAdapterProvider, cfg ParserConfig) *NewsParser {
	return &NewsParser{
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
		cfg:      cfg,
		newsRepo: repo,
	}
}

func (p *NewsParser) Start(ctx context.Context) {
	ticker := time.NewTicker(p.cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("запуск парсинга новостей...")

			news, err := p.ParseActiveChannels(ctx)
			if err != nil {
				log.Printf("ошибка парсинга: %v", err)
				continue
			}
			if len(news) == 0 {
				continue
			}

			inserted, err := p.newsRepo.AddNewsList(news)
			if err != nil {
				log.Printf("ошибка сохранения %d новостей: %v", len(news), err)
				continue
			}
			log.Printf("добавлено %d из %d новостей (%.1f%%)", len(inserted), len(news), 100*float64(len(inserted))/float64(len(news)))

		case <-ctx.Done():
			log.Println("парсер остановлен")
			return
		}
	}
}
