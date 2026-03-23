package v1

import (
	"context"
	"fmt"
	"net/http"

	newsAdapter "github.com/I-Van-Radkov/vesta-gkh/internal/adapter/news"
	"github.com/I-Van-Radkov/vesta-gkh/internal/config"
	"github.com/I-Van-Radkov/vesta-gkh/internal/controller/http/v1/handlers"
	newsUsecase "github.com/I-Van-Radkov/vesta-gkh/internal/usecase/news"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	srv *http.Server
}

func NewServer(cfg *config.Config) *Server {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%v", cfg.Port),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		Handler:      nil,
	}

	return &Server{
		srv: srv,
	}
}

func (s *Server) RegisterHandlers(ctx context.Context, cfg *config.Config, db *pgxpool.Pool) error {
	newsRepo := newsAdapter.NewNewsRepo(db)
	newsUsecase := newsUsecase.NewNewsUsecase(ctx, newsRepo, cfg.ParserConfig)
	newsHandlers := handlers.NewNewsHandlers(newsUsecase)

	router := gin.Default()
	router.GET("/news", newsHandlers.GetNews)

	s.srv.Handler = router

	return nil
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
