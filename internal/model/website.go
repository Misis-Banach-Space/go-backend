package model

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type Website struct {
	Id        uint
	Url       string
	Category  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type WebsiteCreate struct {
	Url string `json:"url" validate:"required,http_url"`
}

type WebsiteDto struct {
	Id       uint   `json:"id"`
	Url      string `json:"url"`
	Category string `json:"category"`
}

type WebsiteCategoryCount struct {
	Category string `json:"category"`
	Count    uint   `json:"count"`
}

type WebsiteRepository interface {
	Add(c *fiber.Ctx, websiteData WebsiteCreate) error
	GetById(c *fiber.Ctx, id uint) (*WebsiteDto, error)
	GetAll(c *fiber.Ctx) (*[]WebsiteDto, error)
	GetByCategory(c *fiber.Ctx, category string) (*[]WebsiteDto, error)
	GetWebsitesCategoryCount(c *fiber.Ctx) (*[]WebsiteCategoryCount, error)
	UpdateCategory(c *fiber.Ctx, websiteId uint, category string) error
}
