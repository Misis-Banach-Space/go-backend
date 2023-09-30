package controller

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/yogenyslav/kokoc-hack/internal/model"
	"github.com/yogenyslav/kokoc-hack/internal/service"
	"github.com/yogenyslav/kokoc-hack/internal/utils"
)

type websiteController struct {
	repository model.WebsiteRepository
	validator  *validator.Validate
	rabbitmq   *service.RabbitMQ
}

func NewWebsiteController(wr model.WebsiteRepository, rabbitmq *service.RabbitMQ) websiteController {
	return websiteController{
		repository: wr,
		validator:  validator.New(validator.WithRequiredStructEnabled()),
		rabbitmq:   rabbitmq,
	}
}

func (wc *websiteController) CreateWebsite(c *fiber.Ctx) error {
	websiteData := model.WebsiteCreate{}
	if err := json.Unmarshal(c.Body(), &websiteData); err != nil {
		return utils.ErrValidationError("websiteData", err)
	}

	if err := wc.validator.Struct(&websiteData); err != nil {
		return utils.ErrValidationError("websiteData", err)
	}

	websiteId, err := wc.repository.Add(c, websiteData)
	if err != nil {
		return utils.ErrCreateRecordsFailed("website", err)
	}

	go wc.rabbitmq.PublishUrl(c.Context(), "url_queue", model.UrlRequest{
		Id:  websiteId,
		Url: websiteData.Url,
	}, wc.repository)

	return c.Status(http.StatusCreated).JSON(websiteId)
}

func (wc *websiteController) GetWebsiteById(c *fiber.Ctx) error {
	websiteId, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrValidationError("id", err)
	}
	if websiteId <= 0 {
		return utils.ErrValidationError("id", errors.New(fmt.Sprintf("id must be positive: %d", websiteId)))
	}

	website, err := wc.repository.GetOneByFilter(c, "id", uint(websiteId))
	if err != nil {
		return utils.ErrGetRecordsFailed("website", err)
	}

	return c.Status(http.StatusOK).JSON(website)
}

func (wc *websiteController) GetWebsiteByUrl(c *fiber.Ctx) error {
	websiteData := model.GetWebsiteByUrlRequest{}
	if err := json.Unmarshal(c.Body(), &websiteData); err != nil {
		return utils.ErrValidationError("websiteUrl", err)
	}

	website, err := wc.repository.GetOneByFilter(c, "url", websiteData.Url)
	if err != nil {
		return utils.ErrGetRecordsFailed("website", err)
	}

	if website.Category == "unmatched" || website.Theme == "unmatched" || website.Stats == nil {
		go wc.rabbitmq.PublishUrl(c.Context(), "url_queue", model.UrlRequest{
			Id:  website.Id,
			Url: websiteData.Url,
		}, wc.repository)
	}

	return c.Status(http.StatusOK).JSON(website)
}

func (wc *websiteController) GetWebsitesByCategory(c *fiber.Ctx) error {
	category := c.Params("category")

	websites, err := wc.repository.GetManyByFilter(c, "category", category)
	if err != nil {
		return utils.ErrGetRecordsFailed("website", err)
	}

	return c.Status(http.StatusOK).JSON(websites)
}

func (wc *websiteController) GetAllWebsites(c *fiber.Ctx) error {
	websites, err := wc.repository.GetManyByFilter(c, "", "")
	if err != nil {
		return utils.ErrGetRecordsFailed("website", err)
	}

	return c.Status(http.StatusOK).JSON(websites)
}

func (wc *websiteController) GetWebsitesCategoryCount(c *fiber.Ctx) error {
	websitesCount, err := wc.repository.GetWebsitesCategoryCount(c)
	if err != nil {
		return utils.ErrGetRecordsFailed("website", err)
	}

	return c.Status(http.StatusOK).JSON(websitesCount)
}

func (wc *websiteController) SseUpdateCategory(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		var i int
		for event := range wc.rabbitmq.Events() {
			i++
			msg := fmt.Sprintf("%d - the event is %s", i, event)
			fmt.Fprintf(w, "data: Message: %s\n\n", msg)

			err := w.Flush()
			if err != nil {
				fmt.Printf("Error while flushing: %v. Closing http connection.\n", err)
				break
			}
		}
	})

	return nil
}
