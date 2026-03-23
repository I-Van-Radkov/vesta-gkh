package handlers

import (
	"net/http"

	"github.com/I-Van-Radkov/vesta-gkh/internal/dto"
	"github.com/gin-gonic/gin"
)

type NewsUsecaseProvider interface {
	GetNewsList() (*dto.NewsListResponse, error)
}

type NewsHandlers struct {
	newsUsecase NewsUsecaseProvider
}

func NewNewsHandlers(usecase NewsUsecaseProvider) *NewsHandlers {
	return &NewsHandlers{
		newsUsecase: usecase,
	}
}

func (h *NewsHandlers) GetNews(c *gin.Context) {
	news, err := h.newsUsecase.GetNewsList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, news)
}
