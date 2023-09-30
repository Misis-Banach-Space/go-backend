package repository

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yogenyslav/kokoc-hack/internal/logging"
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
			theme text default 'unmatched',
			stats jsonb,
			created_at timestamp default current_timestamp,
			updated_at timestamp default current_timestamp
		);
	`)
	return &websiteRepository{
		tableName: tableName,
	}, err
}

func (wr *websiteRepository) Add(c *fiber.Ctx, websiteData model.WebsiteCreate) (uint, error) {
	db := c.Locals("session").(*pgxpool.Conn)

	var websiteId uint
	err := db.QueryRow(c.Context(), `
		insert into `+wr.tableName+`(url)
		values($1) returning "id"
	`, websiteData.Url).Scan(&websiteId)
	if err != nil {
		return 0, err
	}

	return websiteId, nil
}

func (wr *websiteRepository) GetOneByFilter(c *fiber.Ctx, filter string, value any) (*model.WebsiteDto, error) {
	db := c.Locals("session").(*pgxpool.Conn)
	website := &model.WebsiteDto{}

	row := db.QueryRow(c.Context(), `
		select id, url, category, theme, stats from `+wr.tableName+` where `+filter+` = $1 
	`, value)
	err := row.Scan(&website.Id, &website.Url, &website.Category, &website.Theme, &website.Stats)
	if err != nil {
		return nil, err
	}

	return website, nil
}

func (wr *websiteRepository) GetManyByFilter(c *fiber.Ctx, filter string, value any) (*[]model.WebsiteDto, error) {
	db := c.Locals("session").(*pgxpool.Conn)
	var websites []model.WebsiteDto

	var rows pgx.Rows
	var err error
	sql := "select id, url, category, theme, stats from " + wr.tableName
	if filter != "" {
		sql += " where" + filter + " = $1;"
		rows, err = db.Query(c.Context(), sql, value)
	} else {
		rows, err = db.Query(c.Context(), sql)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		website := model.WebsiteDto{}
		err := rows.Scan(&website.Id, &website.Url, &website.Category, &website.Theme, &website.Stats)
		if err != nil {
			return nil, err
		}

		websites = append(websites, website)
	}

	return &websites, nil
}

func (wr *websiteRepository) Update(c context.Context, db *pgxpool.Pool, newData model.UrlResponse) error {
	_, err := db.Exec(c, `
		update `+wr.tableName+` set category = $1 where id = $2;
	`, newData.Category, newData.Id)
	if err != nil {
		logging.Log.Error("category error")
		return err
	}
	_, err = db.Exec(c, `
		update `+wr.tableName+` set theme = $1 where id = $2;
	`, newData.Theme, newData.Id)
	if err != nil {
		logging.Log.Error("theme error")
		return err
	}
	_, err = db.Exec(c, `
		update `+wr.tableName+` set stats = $1 where id = $2;
	`, newData.Stats, newData.Id)
	if err != nil {
		logging.Log.Error("stats error")
		return err
	}

	return nil
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
