package infos

import (
	"github.com/gofiber/fiber/v2"
)

// List godoc
// @Summary      List infos
// @Description  Get all infos
// @Tags         infos
// @Accept       json
// @Produce      json
// @Success      200  {object}  InfoShowModel
// @Router       /infos [get]
// @security     BasicAuth
func (h *handler) List(ctx *fiber.Ctx) error {
	showInfo := InfoShowModel{
		Location:    h.config.Service.Location,
		MaxSyncSize: h.config.Service.MaxSyncSize,
		Message:     h.config.Service.Message,
		Status:      int(h.config.Service.Status),
		Version:     "1.1.13",
	}

	return ctx.
		Status(fiber.StatusOK).
		JSON(showInfo)
}
