package repository

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yogenyslav/kokoc-hack/internal/model"
)

type websiteRepository struct {
	tableName string
}

func NewWebsiteRepository(ctx context.Context, tableName string, db *pgxpool.Conn) (model.WebsiteRepository, error) {
	_, err := db.Exec(ctx, `
		create table if not exists `+tableName+`(
			id serial primary key,
			url text unique,
			category text default 'unmatched',
			created_at timestamp default current_timestamp,
			updated_at timestamp default current_timestamp
		);
	`)
	return &websiteRepository{
		tableName: tableName,
	}, err
}

func (wr *websiteRepository) Add(c *fiber.Ctx, websiteData model.WebsiteCreate) error {
	db := c.Locals("session").(*pgxpool.Conn)

	_, err := db.Exec(c.Context(), `
		insert into `+wr.tableName+`(url)
		values($1) returning "id"
	`, websiteData.Url)

	return err
}

func (wr *websiteRepository) GetById(c *fiber.Ctx, id uint) (*model.WebsiteDto, error) {
	db := c.Locals("session").(*pgxpool.Conn)
	website := &model.WebsiteDto{}

	row := db.QueryRow(c.Context(), `
		select id, url, category from `+wr.tableName+` where id = $1 
	`, id)
	err := row.Scan(&website.Id, &website.Url, &website.Category)
	return website, err
}

func (wr *websiteRepository) GetAll(c *fiber.Ctx) (*[]model.WebsiteDto, error) {
	db := c.Locals("session").(*pgxpool.Conn)
	var websites []model.WebsiteDto

	rows, err := db.Query(c.Context(), `
		select id, url, category from `+wr.tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		website := model.WebsiteDto{}
		err := rows.Scan(&website.Id, &website.Url, &website.Category)
		if err != nil {
			return nil, err
		}

		websites = append(websites, website)
	}

	return &websites, err
}

func (wr *websiteRepository) GetByCategory(c *fiber.Ctx, category string) (*[]model.WebsiteDto, error) {
	allWebsites, err := wr.GetAll(c)
	if err != nil {
		return nil, err
	}

	var websites []model.WebsiteDto
	for _, website := range *allWebsites {
		if website.Category == category {
			websites = append(websites, website)
		}
	}

	return &websites, nil
}

func (wr *websiteRepository) UpdateCategory(c *fiber.Ctx, websiteId uint, category string) error {
	db := c.Locals("session").(*pgxpool.Conn)

	_, err := db.Exec(c.Context(), `
		update `+wr.tableName+` set category = $1 where id = $2
	`, category, websiteId)

	return err
}

func (wr *websiteRepository) GetWebsitesCategoryCount(c *fiber.Ctx) (*[]model.WebsiteCategoryCount, error) {
	db := c.Locals("session").(*pgxpool.Conn)

	rows, err := db.Query(c.Context(), `
		select category, count(*)
		from `+wr.tableName+` group by category
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var websitesCategoryCounts []model.WebsiteCategoryCount
	for rows.Next() {
		cur := model.WebsiteCategoryCount{}
		err := rows.Scan(&cur.Category, &cur.Count)
		if err != nil {
			return nil, err
		}

		websitesCategoryCounts = append(websitesCategoryCounts, cur)
	}

	return &websitesCategoryCounts, nil
}
