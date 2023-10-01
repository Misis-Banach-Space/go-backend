package controller

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/yogenyslav/kokoc-hack/internal/model"
	"github.com/yogenyslav/kokoc-hack/internal/service"
	"github.com/yogenyslav/kokoc-hack/internal/utils"
)

type reportController struct {
	repository model.ReportRepository
	validator  *validator.Validate
	rabbitmq   *service.RabbitMQ
}

func NewReportController(rr model.ReportRepository, rabbitmq *service.RabbitMQ) reportController {
	return reportController{
		repository: rr,
		validator:  validator.New(validator.WithRequiredStructEnabled()),
		rabbitmq:   rabbitmq,
	}
}

func (rc *reportController) GetReport(c *fiber.Ctx) error {
	var wg sync.WaitGroup
	reportRequest := model.ReportRequest{}

	if err := json.Unmarshal(c.Body(), &reportRequest); err != nil {
		return utils.ErrValidationError("urlReportRequest", err)
	}

	if err := rc.validator.Struct(&reportRequest); err != nil {
		return utils.ErrValidationError("urlReportRequest", err)
	}

	urlDomain, err := utils.GetUrlDomain(reportRequest.Url)
	if err != nil {
		return utils.ErrCustomResponse(http.StatusInternalServerError, "failed to get url domain", err)
	}

	var reportResponse model.ReportResponse
	if reportRequest.Url == urlDomain {
		website, err := rc.repository.GetWebsite(c, reportRequest.Url)
		if err != nil {
			return utils.ErrCreateRecordsFailed("website", err)
		}

		wg.Add(1)
		go rc.rabbitmq.PublishUrl(c.Context(), "url_queue", model.UrlRequest{Id: website.Id, Url: website.Url}, rc.repository.GetWebsiteRepository(), &wg)

		reportResponse = model.ReportResponse{
			Category: website.Category,
			Theme:    website.Theme,
		}
	} else {
		page, err := rc.repository.GetPage(c, reportRequest.Url)
		if err != nil {
			return utils.ErrCreateRecordsFailed("page", err)
		}

		wg.Add(1)
		go rc.rabbitmq.PublishUrl(c.Context(), "url_queue", model.UrlRequest{Id: page.Id, Url: page.Url}, rc.repository.GetPageRepository(), &wg)

		reportResponse = model.ReportResponse{
			Category: page.Category,
			Theme:    page.Theme,
		}
	}

	wg.Wait()

	return c.Status(http.StatusOK).JSON(reportResponse)
}
