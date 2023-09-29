package model

import (
	"time"

	"github.com/gofiber/fiber/v2"
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
	Id       uint   `json:"id"`
	Url      string `json:"url"`
	Category string `json:"category"`
	Theme    string `json:"theme"`
}

type PageRepository interface {
	Add(c *fiber.Ctx, pageData PageCreate, websiteId uint) error
	GetById(c *fiber.Ctx, id uint) (*PageDto, error)
	GetPagesByWebsiteId(c *fiber.Ctx, websiteId uint) (*[]PageDto, error)
}
