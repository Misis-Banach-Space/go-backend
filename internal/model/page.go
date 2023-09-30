package model

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Page struct {
	Id          uint
	Url         string
	Category    string
	Theme       string
	FkWebsiteId uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type PageCreate struct {
	Url string `json:"url" validate:"required,http_url"`
}

type PageDto struct {
	Id          uint   `json:"id"`
	Url         string `json:"url"`
	Category    string `json:"category"`
	Theme       string `json:"theme"`
	FkWebsiteId uint   `json:"websiteId"`
}

type GetPageByUrlRequest struct {
	Url string `json:"url"`
}

type PageRepository interface {
	Add(c *fiber.Ctx, pageData PageCreate, websiteId uint) (uint, error)
	GetOneByFilter(c *fiber.Ctx, filter string, value any) (*PageDto, error)
	GetPagesByWebsiteId(c *fiber.Ctx, websiteId uint) (*[]PageDto, error)
	Update(c context.Context, db *pgxpool.Pool, newData UrlResponse) error
}
