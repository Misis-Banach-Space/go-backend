package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/yogenyslav/kokoc-hack/internal/model"
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
		return errValidationError("websiteData", err)
	}

	if err := wc.validator.Struct(&websiteData); err != nil {
		return errValidationError("websiteData", err)
	}

	err := wc.repository.Add(c, websiteData)
	if err != nil {
		return errCreateRecordsFailed("website", err)
	}

	return c.SendStatus(http.StatusCreated)
}

func (wc *websiteController) GetWebsiteById(c *fiber.Ctx) error {
	websiteId, err := c.ParamsInt("id")
	if err != nil {
		return errValidationError("id", err)
	}
	if websiteId <= 0 {
		return errValidationError("id", errors.New(fmt.Sprintf("id must be positive: %d", websiteId)))
	}

	website, err := wc.repository.GetById(c, uint(websiteId))
	if err != nil {
		return errGetRecordsFailed("website", err)
	}

	return c.Status(http.StatusOK).JSON(website)
}

func (wc *websiteController) GetWebsitesByCategory(c *fiber.Ctx) error {
	category := c.Params("category")

	websites, err := wc.repository.GetByCategory(c, category)
	if err != nil {
		return errGetRecordsFailed("website", err)
	}

	return c.Status(http.StatusOK).JSON(websites)
}

func (wc *websiteController) GetAllWebsites(c *fiber.Ctx) error {
	websites, err := wc.repository.GetAll(c)
	if err != nil {
		return errGetRecordsFailed("website", err)
	}

	return c.Status(http.StatusOK).JSON(websites)
}

func (wc *websiteController) UpdateCategory(c *fiber.Ctx) error {
	websiteId := c.QueryInt("id")
	category := c.Query("category")

	if websiteId <= 0 {
		return errValidationError("id", errors.New(fmt.Sprintf("id must be positive: %d", websiteId)))
	}

	err := wc.repository.UpdateCategory(c, uint(websiteId), category)
	if err != nil {
		return errUpdateRecordsFailed("website", err)
	}

	return c.SendStatus(http.StatusNoContent)
}

func (wc *websiteController) GetWebsitesCategoryCount(c *fiber.Ctx) error {
	websitesCount, err := wc.repository.GetWebsitesCategoryCount(c)
	if err != nil {
		return errGetRecordsFailed("website", err)
	}

	return c.Status(http.StatusOK).JSON(websitesCount)
}
