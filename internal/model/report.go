package model

import "github.com/gofiber/fiber/v2"

type ReportRequest struct {
	Url string `json:"url" validate:"required,http_url"`
}

type ReportResponse struct {
	Category string `json:"category"`
	Theme    string `json:"theme"`
}

type ReportRepository interface {
	GetWebsite(c *fiber.Ctx, url string) (*WebsiteDto, error) // create if not exists
	GetPage(c *fiber.Ctx, url string) (*PageDto, error)       // create if not exists
	GetWebsiteRepository() WebsiteRepository
	GetPageRepository() PageRepository
}
