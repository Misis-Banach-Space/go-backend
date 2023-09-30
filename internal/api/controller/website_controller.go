package controller

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/yogenyslav/kokoc-hack/internal/model"
	"github.com/yogenyslav/kokoc-hack/internal/utils"
)

type websiteController struct {
	repository model.WebsiteRepository
	validator  *validator.Validate
}

func NewWebsiteController(websiteRepository model.WebsiteRepository) websiteController {
	return websiteController{
		repository: websiteRepository,
		validator:  validator.New(validator.WithRequiredStructEnabled()),
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

	// send request to ml with rabbit
	// get statistics

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

	website, err := wc.repository.GetById(c, uint(websiteId))
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

	website, err := wc.repository.GetByUrl(c, websiteData.Url)
	if err != nil {
		return utils.ErrGetRecordsFailed("website", err)
	}

	return c.Status(http.StatusOK).JSON(website)
}

func (wc *websiteController) GetWebsitesByCategory(c *fiber.Ctx) error {
	category := c.Params("category")

	websites, err := wc.repository.GetByCategory(c, category)
	if err != nil {
		return utils.ErrGetRecordsFailed("website", err)
	}

	return c.Status(http.StatusOK).JSON(websites)
}

func (wc *websiteController) GetAllWebsites(c *fiber.Ctx) error {
	websites, err := wc.repository.GetAll(c)
	if err != nil {
		return utils.ErrGetRecordsFailed("website", err)
	}

	return c.Status(http.StatusOK).JSON(websites)
}

func (wc *websiteController) UpdateCategory(c *fiber.Ctx) error {
	websiteId := c.QueryInt("id")
	category := c.Query("category")

	if websiteId <= 0 {
		return utils.ErrValidationError("id", errors.New(fmt.Sprintf("id must be positive: %d", websiteId)))
	}

	err := wc.repository.UpdateCategory(c, uint(websiteId), category)
	if err != nil {
		return utils.ErrUpdateRecordsFailed("website", err)
	}

	return c.SendStatus(http.StatusNoContent)
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
		fmt.Println("WRITER")
		var i int
		for {
			i++
			msg := fmt.Sprintf("%d - the time is %v", i, time.Now())
			fmt.Fprintf(w, "data: Message: %s\n\n", msg)
			fmt.Println(msg)

			err := w.Flush()
			if err != nil {
				fmt.Printf("Error while flushing: %v. Closing http connection.\n", err)
				break
			}
			time.Sleep(2 * time.Second)
		}
	})

	return nil
}
