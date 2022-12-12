package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mrusme/xbsapi/api/v1/bookmarks"
	"github.com/mrusme/xbsapi/api/v1/infos"
	"github.com/mrusme/xbsapi/lib"
)

func Register(
	xbsctx *lib.XBSContext,
	fiberRouter *fiber.Router,
) {
	v1 := (*fiberRouter).Group("/v1")

	bookmarks.Register(
		xbsctx,
		&v1,
	)

	infos.Register(
		xbsctx,
		&v1,
	)
}
