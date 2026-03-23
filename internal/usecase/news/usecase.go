package news

import (
	"context"
	"fmt"
	"time"

	"github.com/I-Van-Radkov/vesta-gkh/internal/dto"
	"github.com/I-Van-Radkov/vesta-gkh/internal/models"
	"github.com/google/uuid"
)

type NewsAdapterProvider interface {
	GetKeywordsList() ([]*models.GKHKeyword, error)
	GetRegexRulesList() ([]*models.GKHRegexRule, error)
	GetCategoryWithKeywords() ([]*models.GHKCategoryWithKeywords, error)
	GetNegativeRegexList() ([]*models.GKHNegativeRegex, error)

	GetChannelById(id uuid.UUID) (*models.Channel, error)
	GetActiveChannels() ([]*models.Channel, error)
	GetLatestPubDateByChannel(ctx context.Context, channelID uuid.UUID) (time.Time, error)

	AddNewsList(news []*models.News) ([]uuid.UUID, error)
	GetNewsList() ([]*models.News, error)
	GetNewsById(id uuid.UUID) (*models.News, error)
}

type NewsUsecase struct {
	parser *NewsParser

	newsRepo NewsAdapterProvider
}

func NewNewsUsecase(ctx context.Context, repo NewsAdapterProvider, cfg ParserConfig) *NewsUsecase {
	parser := NewNewsParser(repo, cfg)
	np := NewsUsecase{
		parser:   parser,
		newsRepo: repo,
	}

	go parser.Start(ctx)

	return &np
}

func (u *NewsUsecase) GetNewsList() (*dto.NewsListResponse, error) {
	newsModels, err := u.newsRepo.GetNewsList()
	if err != nil {
		return nil, fmt.Errorf("...") // дописать ошибку
	}

	output := &dto.NewsListResponse{
		Total: int64(len(newsModels)),
	}

	for _, model := range newsModels {
		var channelTitle, sourceName string
		if ch, err := u.newsRepo.GetChannelById(model.ChannelID); err == nil && ch != nil {
			channelTitle = ch.Title
			sourceName = ch.SourceName
		} else {
			channelTitle = "Unknown"
			sourceName = "Unknown"
		}

		item := dto.NewsItem{
			ID:             model.ID,
			ChannelID:      model.ChannelID,
			ChannelTitle:   channelTitle,
			ChannelSource:  sourceName,
			Title:          model.Title,
			Link:           model.Link,
			Description:    model.Description,
			PubDate:        model.PubDate,
			FetchedAt:      model.FetchedAt,
			RelevanceScore: model.RelevanceScore,
			ImageURL:       model.EnclosureURL,
		}

		output.Items = append(output.Items, item)
	}

	return output, nil
}
