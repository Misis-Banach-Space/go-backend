package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/yogenyslav/kokoc-hack/internal/logging"
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
	logging.Log.Debug("test")

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

	if string(reportRequest.Url[len(reportRequest.Url)-1]) == "/" {
		reportRequest.Url = string(reportRequest.Url[:len(reportRequest.Url)-1])
	}
	if reportRequest.Url == urlDomain {
		website, err := rc.repository.GetWebsite(c, reportRequest.Url)
		if err != nil {
			return utils.ErrCreateRecordsFailed("website", err)
		}

		logging.Log.Debugf("website url %s", reportRequest.Url)
		if website.Category == "" || website.Theme == "" || website.Stats == nil {
			wg.Add(1)
			go rc.rabbitmq.PublishUrlWithWaitGroup(c.Context(), "url_queue", model.UrlRequest{Id: website.Id, Url: website.Url}, rc.repository.GetWebsiteRepository(), &wg)
		}
	} else {
		page, err := rc.repository.GetPage(c, reportRequest.Url)
		if err != nil {
			return utils.ErrCreateRecordsFailed("page", err)
		}

		logging.Log.Debugf("page url %s", reportRequest.Url)
		if page.Category == "" || page.Theme == "" {
			wg.Add(1)
			go rc.rabbitmq.PublishUrlWithWaitGroup(c.Context(), "url_queue", model.UrlRequest{Id: page.Id, Url: page.Url}, rc.repository.GetPageRepository(), &wg)
		}
	}
	wg.Wait()

	for event := range rc.rabbitmq.Events() {
		if event.Url == reportRequest.Url {
			return c.Status(http.StatusOK).JSON(model.ReportResponse{
				Category: event.Category,
				Theme:    event.Theme,
			})
		}
	}

	return utils.ErrCustomResponse(http.StatusBadRequest, "can't determine category", errors.New("invalid request"))
}
