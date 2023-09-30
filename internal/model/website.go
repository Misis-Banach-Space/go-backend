package model

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Website struct {
	Id        uint
	Url       string
	Category  string
	Theme     string
	Stats     map[string]interface{}
	CreatedAt time.Time
	UpdatedAt time.Time
}

type WebsiteCreate struct {
	Url string `json:"url" validate:"required,http_url"`
}

type WebsiteDto struct {
	Id       uint                   `json:"id"`
	Url      string                 `json:"url"`
	Category string                 `json:"category"`
	Theme    string                 `json:"theme"`
	Stats    map[string]interface{} `json:"stats"`
}

type WebsiteCategoryCount struct {
	Category string `json:"category"`
	Count    uint   `json:"count"`
}

type GetWebsiteByUrlRequest struct {
	Url string `json:"url"`
}

type WebsiteRepository interface {
	Add(c *fiber.Ctx, websiteData WebsiteCreate) (uint, error)
	GetById(c *fiber.Ctx, id uint) (*WebsiteDto, error)
	GetByUrl(c *fiber.Ctx, url string) (*WebsiteDto, error)
	GetAll(c *fiber.Ctx) (*[]WebsiteDto, error)
	GetByCategory(c *fiber.Ctx, category string) (*[]WebsiteDto, error)
	GetByTheme(c *fiber.Ctx, theme string) (*[]WebsiteDto, error)
	GetWebsitesCategoryCount(c *fiber.Ctx) (*[]WebsiteCategoryCount, error)
	Update(c context.Context, db *pgxpool.Pool, id uint, category, theme string, stats map[string]interface{}) error
}
