package model

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UrlRequest struct {
	Id  uint   `json:"id"`
	Url string `json:"url"`
}

type UrlResponse struct {
	Id       uint                   `json:"id"`
	Url      string                 `json:"url"`
	Stats    map[string]interface{} `json:"stats"`
	Category string                 `json:"category"`
	Theme    string                 `json:"theme"`
}

type UrlEventRepository interface {
	Update(c context.Context, db *pgxpool.Pool, newData UrlResponse) error
}
