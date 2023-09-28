package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/yogenyslav/kokoc-hack/internal/database"
)

func DbSessionMiddleware(db *database.Postgres) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.Context(), time.Duration(db.Timeout())*time.Second)
		defer cancel()

		session, err := db.GetSessionWithContext(ctx)
		if err != nil {
			return fiber.NewError(http.StatusInternalServerError, err.Error())
		}
		defer session.Release()

		if err != nil {
			return fiber.NewError(http.StatusInternalServerError, err.Error())
		}

		c.Locals("session", session)
		return c.Next()
	}
}
