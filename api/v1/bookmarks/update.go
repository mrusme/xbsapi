package bookmarks

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/mrusme/xbsapi/ent/bookmark"

	"github.com/gofiber/fiber/v2"
)

type BookmarkUpdateResponse struct {
	Success  bool               `json:"success"`
	Bookmark *BookmarkShowModel `json:"bookmark"`
	Message  string             `json:"message"`
}

// Update godoc
// @Summary      Update a bookmark sync
// @Description  Update an existing bookmark sync
// @Tags         bookmarks
// @Accept       json
// @Produce      json
// @Param        id   path      string true "Bookmark ID"
// @Param        bookmark body  BookmarkUpdateModel true "Update"
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

	if len(updateBookmark.Bookmarks) > h.config.Service.MaxSyncSize {
		return ctx.
			Status(fiber.StatusRequestEntityTooLarge).
			JSON(fiber.Map{
				"success": false,
				"message": "Quota exceeded!",
			})
	}

	lastUpdated, err := time.Parse(LAST_UPDATED_FORMAT, updateBookmark.LastUpdated)
	if err != nil {
		return ctx.
			Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
	}

	now := time.Now().Truncate(time.Millisecond)
	dbBookmarkN, err := h.entClient.Bookmark.
		Update().
		Where(bookmark.And(
			bookmark.ID(id),
			bookmark.LastUpdatedGTE(lastUpdated.Add(time.Second*-1)),
			bookmark.LastUpdatedLTE(lastUpdated.Add(time.Second*1)),
		)).
		SetBookmarks(updateBookmark.Bookmarks).
		SetLastUpdated(now).
		Save(context.Background())

	if err != nil {
		return ctx.
			Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
	}

	if dbBookmarkN != 1 {
		return ctx.
			Status(fiber.StatusConflict).
			JSON(fiber.Map{
				"success": false,
				"message": "Internal sync error",
			})
	}

	showBookmark := BookmarkShowModel{
		LastUpdated: now.Format(LAST_UPDATED_FORMAT),
	}

	return ctx.
		Status(fiber.StatusOK).
		JSON(showBookmark)
}
