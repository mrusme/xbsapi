package bookmarks

import (
	"context"

	"github.com/mrusme/xbsapi/ent/bookmark"

	"github.com/gofiber/fiber/v2"
)

// Show godoc
// @Summary      Show a bookmark sync lastUpdated
// @Description  Get bookmark sync lastUpdated by ID
// @Tags         bookmarks
// @Accept       json
// @Produce      json
// @Param        id   path      string true "Bookmark sync ID"
// @Success      200  {object}  BookmarkShowResponse
// @Failure      400  {object}  BookmarkShowResponse
// @Failure      404  {object}  BookmarkShowResponse
// @Failure      500  {object}  BookmarkShowResponse
// @Router       /bookmarks/{id} [get]
// @security     BasicAuth
func (h *handler) ShowLastUpdated(ctx *fiber.Ctx) error {
	var err error

	id, err := h.getUUIDFromID(ctx.Params("id"))
	if err != nil {
		return ctx.
			Status(fiber.StatusBadRequest).
			JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
	}

	dbBookmark, err := h.entClient.Bookmark.
		Query().
		Where(
			bookmark.ID(id),
		).
		Only(context.Background())
	if err != nil {
		return ctx.
			Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
	}

	showBookmark := BookmarkShowModel{
		LastUpdated: dbBookmark.LastUpdated.Format(LAST_UPDATED_FORMAT),
	}

	return ctx.
		Status(fiber.StatusOK).
		JSON(showBookmark)
}
