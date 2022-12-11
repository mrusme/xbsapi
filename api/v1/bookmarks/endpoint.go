package bookmarks

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mrusme/xbsapi/lib"
	"go.uber.org/zap"
)

type handler struct {
	xbsctx *lib.XBSContext
	config *lib.Config
	// entClient *ent.Client
	logger *zap.Logger
}

type BookmarkShowModel struct {
	ID           string `json:"id"`
	Bookmarkname string `json:"username"`
	Role         string `json:"role"`
}

type BookmarkCreateModel struct {
	Bookmarkname string `json:"username" validate:"required,alphanum,max=32"`
	Password     string `json:"password" validate:"required"`
	Role         string `json:"role" validate:"required"`
}

type BookmarkUpdateModel struct {
	Password string `json:"password,omitempty" validate:"omitempty,min=5"`
	Role     string `json:"role,omitempty" validate:"omitempty"`
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
