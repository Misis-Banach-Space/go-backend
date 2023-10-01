package router

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/yogenyslav/kokoc-hack/internal/api/controller"
	"github.com/yogenyslav/kokoc-hack/internal/repository"
)

func (r *router) setupReportRoutes(a *fiber.App) error {
	ctx := context.Background()

	session, err := r.db.GetSessionWithContext(ctx)
	if err != nil {
		return err
	}

	websiteRepository, err := repository.NewWebsiteRepository(ctx, "websites", session)
	if err != nil {
		return err
	}
	pageRepository, err := repository.NewPageRepository(ctx, "pages", session)
	if err != nil {
		return err
	}
	reportRepository := repository.NewReportRepository(websiteRepository, pageRepository)

	reportController := controller.NewReportController(reportRepository, r.rabbitmq)

	a.Post("/check_url", reportController.GetReport)

	return nil
}
