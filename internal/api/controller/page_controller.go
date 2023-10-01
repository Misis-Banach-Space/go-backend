package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/yogenyslav/kokoc-hack/internal/model"
	"github.com/yogenyslav/kokoc-hack/internal/service"
	"github.com/yogenyslav/kokoc-hack/internal/utils"
)

type pageController struct {
	pageRepository    model.PageRepository
	websiteRepository model.WebsiteRepository
	validator         *validator.Validate
	rabbitmq          *service.RabbitMQ
}

func NewPageController(pr model.PageRepository, wr model.WebsiteRepository, rabbit *service.RabbitMQ) pageController {
	return pageController{
		pageRepository:    pr,
		websiteRepository: wr,
		validator:         validator.New(validator.WithRequiredStructEnabled()),
		rabbitmq:          rabbit,
	}
}

func (pc *pageController) CreatePage(c *fiber.Ctx) error {
	pageData := model.PageCreate{}
	if err := json.Unmarshal(c.Body(), &pageData); err != nil {
		return utils.ErrValidationError("pageData", err)
	}

	if err := pc.validator.Struct(&pageData); err != nil {
		return utils.ErrValidationError("pageData", err)
	}

	pageDomain, err := utils.GetUrlDomain(pageData.Url)
	if err != nil {
		return utils.ErrCustomResponse(http.StatusBadRequest, err.Error(), err)
	}

	var websiteId uint
	website, err := pc.websiteRepository.GetOneByFilter(c, "url", pageDomain)
	if errors.Is(err, pgx.ErrNoRows) {
		websiteId, err = pc.websiteRepository.Add(c, model.WebsiteCreate{Url: pageDomain})
		if err != nil {
			return utils.ErrCreateRecordsFailed("website", err)
		}
		go pc.rabbitmq.PublishUrl(c.Context(), "url_queue", model.UrlRequest{
			Id:  websiteId,
			Url: pageDomain,
		}, pc.websiteRepository)
	} else if err != nil {
		return utils.ErrGetRecordsFailed("websiteId", err)
	} else {
		websiteId = website.Id
	}

	pageId, err := pc.pageRepository.Add(c, pageData, websiteId)
	if err != nil {
		return utils.ErrCreateRecordsFailed("page", err)
	}

	go pc.rabbitmq.PublishUrl(c.Context(), "url_queue", model.UrlRequest{
		Id:  pageId,
		Url: pageData.Url,
	}, pc.pageRepository)

	return c.Status(http.StatusCreated).JSON(pageId)
}

func (pc *pageController) GetPageById(c *fiber.Ctx) error {
	pageId, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrValidationError("id", err)
	}
	if pageId <= 0 {
		return utils.ErrValidationError("id", errors.New(fmt.Sprintf("id must be positive: %d", pageId)))
	}

	page, err := pc.pageRepository.GetOneByFilter(c, "id", uint(pageId))
	if err != nil {
		return utils.ErrGetRecordsFailed("page", err)
	}

	return c.Status(http.StatusOK).JSON(page)
}

func (pc *pageController) GetPageByUrl(c *fiber.Ctx) error {
	pageData := model.GetPageByUrlRequest{}
	if err := json.Unmarshal(c.Body(), &pageData); err != nil {
		return utils.ErrValidationError("pageUrl", err)
	}

	page, err := pc.pageRepository.GetOneByFilter(c, "url", pageData.Url)
	if err != nil {
		return utils.ErrGetRecordsFailed("page", err)
	}

	if page.Category == "" || page.Theme == "" {
		go pc.rabbitmq.PublishUrl(c.Context(), "url_queue", model.UrlRequest{
			Id:  page.Id,
			Url: pageData.Url,
		}, pc.pageRepository)
	}

	return c.Status(http.StatusOK).JSON(page)
}

func (pc *pageController) GetPageByWebsiteId(c *fiber.Ctx) error {
	websiteId, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrValidationError("id", err)
	}
	if websiteId <= 0 {
		return utils.ErrValidationError("id", errors.New(fmt.Sprintf("id must be positive: %d", websiteId)))
	}

	pages, err := pc.pageRepository.GetPagesByWebsiteId(c, uint(websiteId))
	if err != nil {
		return utils.ErrGetRecordsFailed("pages", err)
	}

	return c.Status(http.StatusOK).JSON(pages)
}
