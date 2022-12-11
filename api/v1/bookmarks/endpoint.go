package bookmarks

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mrusme/xbsapi/ent"
	"github.com/mrusme/xbsapi/lib"
	"go.uber.org/zap"
)

type handler struct {
	xbsctx    *lib.XBSContext
	config    *lib.Config
	entClient *ent.Client
	logger    *zap.Logger
}

type BookmarkShowModel struct {
	ID          string `json:"id,omitempty"`
	Bookmarks   string `json:"bookmarks,omitempty"`
	LastUpdated string `json:"lastUpdated,omitempty"`
	Version     string `json:"version,omitempty"`
}

type BookmarkCreateModel struct {
	Version string `json:"version" validate:"required"`
}

type BookmarkUpdateModel struct {
	Bookmarks   string `json:"bookmarks" validate:"required"`
	LastUpdated string `json:"lastUpdated" validate:"required"`
}

func Register(
	xbsctx *lib.XBSContext,
	fiberRouter *fiber.Router,
) {
	endpoint := new(handler)
	endpoint.xbsctx = xbsctx
	endpoint.config = endpoint.xbsctx.Config
	// endpoint.entClient = endpoint.xbsctx.EntClient
	endpoint.logger = endpoint.xbsctx.Logger

	// bookmarksRouter := (*fiberRouter).Group("/bookmarks")
	// bookmarksRouter.Get("/", endpoint.List)
	// bookmarksRouter.Get("/:id", endpoint.Show)
	// bookmarksRouter.Post("/", endpoint.Create)
	// bookmarksRouter.Put("/:id", endpoint.Update)
	// bookmarksRouter.Delete("/:id", endpoint.Destroy)
}
