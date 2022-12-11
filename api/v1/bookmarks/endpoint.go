package bookmarks

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

func (h *handler) getUUIDFromID(param_id string) (uuid.UUID, error) {
	param_uuid := fmt.Sprintf("%s-%s-%s-%s-%s",
		param_id[0:8],
		param_id[8:12],
		param_id[12:16],
		param_id[16:20],
		param_id[20:32],
	)
	return uuid.Parse(param_uuid)
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
	endpoint.entClient = endpoint.xbsctx.EntClient
	endpoint.logger = endpoint.xbsctx.Logger

	bookmarksRouter := (*fiberRouter).Group("/bookmarks")

	bookmarksRouter.Get("/:id", endpoint.Show)
	bookmarksRouter.Get("/:id/lastUpdated", endpoint.ShowLastUpdated)
	bookmarksRouter.Get("/:id/version", endpoint.ShowVersion)
	bookmarksRouter.Post("/", endpoint.Create)
	bookmarksRouter.Put("/:id", endpoint.Update)
}
