package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/yogenyslav/kokoc-hack/internal/api/middleware"
	"github.com/yogenyslav/kokoc-hack/internal/database"
)

type router struct {
	db     *database.Postgres
	engine *fiber.App
}

func NewRouter(db *database.Postgres) router {
	app := fiber.New()

	return router{
		db:     db,
		engine: app,
	}
}

func (r *router) Run(addr string) error {
	if err := r.Setup(); err != nil {
		return err
	}

	return r.engine.Listen(addr)
}

func (r *router) Setup() error {
	r.engine.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "*",
		AllowHeaders: "*",
	}))
	r.engine.Use(logger.New(logger.ConfigDefault))
	r.engine.Use(recover.New(recover.ConfigDefault))
	r.engine.Use(middleware.DbSessionMiddleware(r.db))

	apiV1 := r.engine.Group("/api/v1")
	if err := r.setupWebsiteRoutes(apiV1); err != nil {
		return err
	}

	// if err := r.setupWebsiteRoutes(apiV1); err != nil {
	// 	return err
	// }

	// if err := r.db.CreateJoinTables(); err != nil {
	// 	return err
	// }

	return nil
}
