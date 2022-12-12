package bookmarks

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/gofiber/fiber/v2"
)

type BookmarkUpdateResponse struct {
	Success  bool               `json:"success"`
	Bookmark *BookmarkShowModel `json:"bookmark"`
	Message  string             `json:"message"`
}

// Update godoc
// @Summary      Update a bookmark
// @Description  Change an existing bookmark
// @Tags         bookmarks
// @Accept       json
// @Produce      json
// @Param        id   path      string true "Bookmark ID"
// @Param        bookmark body      BookmarkUpdateModel true "Change bookmark"
// @Success      200  {object}  BookmarkUpdateResponse
// @Failure      400  {object}  BookmarkUpdateResponse
// @Failure      404  {object}  BookmarkUpdateResponse
// @Failure      500  {object}  BookmarkUpdateResponse
// @Router       /bookmarks/{id} [put]
// @security     BasicAuth
func (h *handler) Update(ctx *fiber.Ctx) error {
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

	updateBookmark := new(BookmarkUpdateModel)
	if err = ctx.BodyParser(updateBookmark); err != nil {
		return ctx.
			Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
	}

	validate := validator.New()
	if err = validate.Struct(*updateBookmark); err != nil {
		return ctx.
			Status(fiber.StatusBadRequest).
			JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
	}

	dbBookmarkTmp := h.entClient.Bookmark.
		UpdateOneID(id)

	if updateBookmark.Bookmarks != "" {
		dbBookmarkTmp = dbBookmarkTmp.
			SetBookmarks(updateBookmark.Bookmarks)
	}

	if updateBookmark.LastUpdated != "" {
		t, err := time.Parse(LAST_UPDATED_FORMAT, updateBookmark.LastUpdated)
		if err != nil {
			return ctx.
				Status(fiber.StatusInternalServerError).
				JSON(fiber.Map{
					"success": false,
					"message": err.Error(),
				})
		}
		dbBookmarkTmp = dbBookmarkTmp.
			SetLastUpdated(t)
	}

	dbBookmark, err := dbBookmarkTmp.Save(context.Background())

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
