package news

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type NewsRepo struct {
	db *pgxpool.Pool
}

func NewNewsRepo(db *pgxpool.Pool) *NewsRepo {
	return &NewsRepo{db: db}
}
