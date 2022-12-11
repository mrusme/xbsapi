package bookmarks

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/mrusme/xbsapi/ent/bookmark"

	"github.com/gofiber/fiber/v2"
)

type BookmarkShowResponse struct {
	Success  bool               `json:"success"`
	Bookmark *BookmarkShowModel `json:"bookmark"`
	Message  string             `json:"message"`
}

// Show godoc
// @Summary      Show a bookmark
// @Description  Get bookmark by ID
// @Tags         bookmarks
// @Accept       json
// @Produce      json
// @Param        id   path      string true "Bookmark ID"
// @Success      200  {object}  BookmarkShowResponse
// @Failure      400  {object}  BookmarkShowResponse
// @Failure      404  {object}  BookmarkShowResponse
// @Failure      500  {object}  BookmarkShowResponse
// @Router       /bookmarks/{id} [get]
// @security     BasicAuth
func (h *handler) Show(ctx *fiber.Ctx) error {
	var err error

	param_id := ctx.Params("id")
	param_uuid := fmt.Sprintf("%s-%s-%s-%s-%s",
		param_id[0:8],
		param_id[8:12],
		param_id[12:16],
		param_id[16:20],
		param_id[20:32],
	)
	fmt.Println(param_uuid)
	id, err := uuid.Parse(param_uuid)
	if err != nil {
		return ctx.
			Status(fiber.StatusBadRequest).
			JSON(BookmarkShowResponse{
				Success:  false,
				Bookmark: nil,
				Message:  err.Error(),
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
		Bookmarks:   dbBookmark.Bookmarks,
		LastUpdated: dbBookmark.LastUpdated.Format(LAST_UPDATED_FORMAT),
		Version:     dbBookmark.Version,
	}

	return ctx.
		Status(fiber.StatusOK).
		JSON(showBookmark)
}
