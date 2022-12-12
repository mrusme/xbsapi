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
// @Success      200  {object}  InfoListResponse
// @Failure      400  {object}  InfoListResponse
// @Failure      404  {object}  InfoListResponse
// @Failure      500  {object}  InfoListResponse
// @Router       /infos [get]
// @security     BasicAuth
func (h *handler) List(ctx *fiber.Ctx) error {
	showInfo := InfoShowModel{
		MaxSyncSize: 204800,
		Message:     "It really whips the llama's ass",
		Status:      1,
		Version:     "2.0.0",
	}

	return ctx.
		Status(fiber.StatusOK).
		JSON(showInfo)
}
