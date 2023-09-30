package router

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/yogenyslav/kokoc-hack/internal/api/controller"
	"github.com/yogenyslav/kokoc-hack/internal/repository"
)

func (r *router) setupPageRoutes(g fiber.Router) error {
	ctx := context.Background()

	session, err := r.db.GetSessionWithContext(ctx)
	if err != nil {
		return err
	}

	pageRepository, err := repository.NewPageRepository(ctx, "pages", session)
	if err != nil {
		return err
	}
	websiteRepository, err := repository.NewWebsiteRepository(ctx, "websites", session)
	pageController := controller.NewPageController(pageRepository, websiteRepository)

	pages := g.Group("/pages")
	pages.Get("/by_website/:id", pageController.GetPageByWebsiteId)
	pages.Get("/:id", pageController.GetPageById)
	pages.Post("/create", pageController.CreatePage)
	pages.Post("/check_url", pageController.GetPageByUrl)

	return nil
}
