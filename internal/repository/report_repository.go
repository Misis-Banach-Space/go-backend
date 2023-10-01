package repository

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/yogenyslav/kokoc-hack/internal/model"
)

type reportRepository struct {
	wr model.WebsiteRepository
	pr model.PageRepository
}

func NewReportRepository(wr model.WebsiteRepository, pr model.PageRepository) model.ReportRepository {
	return &reportRepository{
		wr: wr,
		pr: pr,
	}
}

func (rr *reportRepository) GetWebsite(c *fiber.Ctx, url string) (*model.WebsiteDto, error) {
	website, err := rr.wr.GetOneByFilter(c, "url", url)
	if errors.Is(err, pgx.ErrNoRows) {
		websiteId, err := rr.wr.Add(c, model.WebsiteCreate{Url: url})
		if err != nil {
			return nil, err
		}

		website, err = rr.wr.GetOneByFilter(c, "id", websiteId)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	return website, nil
}
func (rr *reportRepository) GetPage(c *fiber.Ctx, url string) (*model.PageDto, error) {
	page, err := rr.pr.GetOneByFilter(c, "url", url)
	if errors.Is(err, pgx.ErrNoRows) {
		website, err := rr.GetWebsite(c, url)
		if err != nil {
			return nil, err
		}

		pageId, err := rr.pr.Add(c, model.PageCreate{Url: url}, uint(website.Id))
		if err != nil {
			return nil, err
		}

		page, err = rr.pr.GetOneByFilter(c, "id", pageId)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	return page, nil
}

func (rr *reportRepository) GetWebsiteRepository() model.WebsiteRepository {
	return rr.wr
}

func (rr *reportRepository) GetPageRepository() model.PageRepository {
	return rr.pr
}
