package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	v1.Register(
		xbsctx,
		&api,
	)
}
