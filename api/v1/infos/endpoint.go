package infos

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

type InfoShowModel struct {
	Location    string `json:"location"`
	MaxSyncSize int    `json:"maxSyncSize"`
	Message     string `json:"message"`
	Status      int    `json:"status"`
	Version     string `json:"version"`
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

	infosRouter := (*fiberRouter).Group("/info")

	infosRouter.Get("/", endpoint.List)
}
