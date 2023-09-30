package router

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/yogenyslav/kokoc-hack/internal/api/controller"
	"github.com/yogenyslav/kokoc-hack/internal/repository"
)

func (r *router) setupWebsiteRoutes(g fiber.Router) error {
	ctx := context.Background()

	session, err := r.db.GetSessionWithContext(ctx)
	if err != nil {
		return err
	}

	websiteRepository, err := repository.NewWebsiteRepository(ctx, "websites", session)
	if err != nil {
		return err
	}

	websiteController := controller.NewWebsiteController(websiteRepository, r.rabbitmq)

	websites := g.Group("/websites")
	websites.Get("/all", websiteController.GetAllWebsites)
	websites.Get("/count", websiteController.GetWebsitesCategoryCount)
	websites.Get("/category/:category", websiteController.GetWebsitesByCategory)
	websites.Get("/:id", websiteController.GetWebsiteById)
	websites.Post("/create", websiteController.CreateWebsite)
	websites.Post("/check_url", websiteController.GetWebsiteByUrl)

	websites.Get("/sse_update", websiteController.SseUpdateCategory)

	return nil
}
