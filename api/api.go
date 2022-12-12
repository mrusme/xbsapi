package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	v1 "github.com/mrusme/xbsapi/api/v1"
	"github.com/mrusme/xbsapi/lib"
)

// @title xBrowserSync API
// @version 1.0
// @description The xBrowserSync REST API v1

// @contact.name Marius
// @contact.url https://xn--gckvb8fzb.com
// @contact.email marius@xn--gckvb8fzb.com

// @license.name GPL-3.0
// @license.url https://github.com/mrusme/xbsapi/blob/master/LICENSE

// @host localhost:8000
// @BasePath /api/v1
// @accept json
// @produce json
// @schemes http
// @securityDefinitions.basic  BasicAuth
func Register(
	xbsctx *lib.XBSContext,
	fiberApp *fiber.App,
) {
	api := fiberApp.Group("/api")
	api.Use(cors.New())
	api.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == "127.0.0.1" && c.Get("x-forwarded-for") == ""
		},
		Max:        xbsctx.Config.Limiter.Max,
		Expiration: xbsctx.Config.Limiter.Expiration,
		KeyGenerator: func(c *fiber.Ctx) string {
			return fmt.Sprintf(
				"%s-%s",
				c.IP(),
				c.Get("x-forwarded-for"),
			)
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.
				Status(fiber.StatusTooManyRequests).
				JSON(fiber.Map{})
		},
	}))

	v1.Register(
		xbsctx,
		&api,
	)
}
