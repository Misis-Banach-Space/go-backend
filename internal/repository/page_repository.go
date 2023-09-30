package repository

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yogenyslav/kokoc-hack/internal/model"
)

type pageRepository struct {
	tableName string
}

func NewPageRepository(ctx context.Context, tableName string, db *pgxpool.Conn) (model.PageRepository, error) {
	_, err := db.Exec(ctx, `
		create table if not exists `+tableName+`(
			id serial primary key,
			url text unique,
			category text default 'unmatched',
			theme text default 'unmatched',
			fk_website_id int not null,
			created_at timestamp default current_timestamp,
			updated_at timestamp default current_timestamp,
			foreign key(fk_website_id) references websites(id)
		);
	`)
	return &pageRepository{
		tableName: tableName,
	}, err
}

func (pr *pageRepository) Add(c *fiber.Ctx, pageData model.PageCreate, websiteId uint) (uint, error) {
	db := c.Locals("session").(*pgxpool.Conn)

	var pageId uint
	err := db.QueryRow(c.Context(), `
		insert into `+pr.tableName+`(url, fk_website_id)
		values($1, $2) returning "id"
	`, pageData.Url, websiteId).Scan(&pageId)
	if err != nil {
		return 0, err
	}

	return pageId, nil
}

func (pr *pageRepository) GetById(c *fiber.Ctx, id uint) (*model.PageDto, error) {
	db := c.Locals("session").(*pgxpool.Conn)
	page := &model.PageDto{}

	row := db.QueryRow(c.Context(), `
		select id, url, category, theme from `+pr.tableName+` where id = $1 
	`, id)
	err := row.Scan(&page.Id, &page.Url, &page.Category, &page.Theme)
	if err != nil {
		return nil, err
	}

	return page, nil
}

func (pr *pageRepository) GetPagesByWebsiteId(c *fiber.Ctx, websiteId uint) (*[]model.PageDto, error) {
	db := c.Locals("session").(*pgxpool.Conn)
	var pages []model.PageDto

	rows, err := db.Query(c.Context(), `
		select id, url, category, theme from `+pr.tableName+` where fk_website_id = $1
	`, websiteId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		page := model.PageDto{}
		err := rows.Scan(&page.Id, &page.Url, &page.Category, &page.Theme)
		if err != nil {
			return nil, err
		}

		pages = append(pages, page)
	}

	return &pages, nil
}

func (pr *pageRepository) GetWebsiteIdByDomain(c *fiber.Ctx, domain string) (uint, error) {
	db := c.Locals("session").(*pgxpool.Conn)

	var websiteId uint
	err := db.QueryRow(c.Context(), `
		select id from websites where url = $1
	`, domain).Scan(&websiteId)
	if err != nil {
		return 0, err
	}

	return websiteId, nil
}
