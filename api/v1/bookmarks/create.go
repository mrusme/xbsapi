package bookmarks

import (
	"context"
	// "github.com/google/uuid"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"github.com/gofiber/fiber/v2"
)

type BookmarkCreateResponse struct {
	Success  bool               `json:"success"`
	Bookmark *BookmarkShowModel `json:"bookmark"`
	Message  string             `json:"message"`
}

// Create godoc
// @Summary      Create a bookmark
// @Description  Add a new bookmark
// @Tags         bookmarks
// @Accept       json
// @Produce      json
// @Param        bookmark body      BookmarkCreateModel true "Add bookmark"
// @Success      200  {object}  BookmarkCreateResponse
// @Failure      400  {object}  BookmarkCreateResponse
// @Failure      404  {object}  BookmarkCreateResponse
// @Failure      500  {object}  BookmarkCreateResponse
// @Router       /bookmarks [post]
// @security     BasicAuth
func (h *handler) Create(ctx *fiber.Ctx) error {
	var err error

	createBookmark := new(BookmarkCreateModel)
	if err = ctx.BodyParser(createBookmark); err != nil {
		h.logger.Debug(
			"Body parsing failed",
			zap.Error(err),
		)
		return ctx.
			Status(fiber.StatusInternalServerError).
			JSON(BookmarkCreateResponse{
				Success:  false,
				Bookmark: nil,
				Message:  err.Error(),
			})
	}

	validate := validator.New()
	if err = validate.Struct(*createBookmark); err != nil {
		h.logger.Debug(
			"Validation failed",
			zap.Error(err),
		)
		return ctx.
			Status(fiber.StatusBadRequest).
			JSON(BookmarkCreateResponse{
				Success:  false,
				Bookmark: nil,
				Message:  err.Error(),
			})
	}

	dbBookmark, err := h.entClient.Bookmark.
		Create().
		SetBookmarks("").
		SetVersion(createBookmark.Version).
		Save(context.Background())

	if err != nil {
		h.logger.Debug(
			"Could not create bookmark",
			zap.Error(err),
		)
		return ctx.
			Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
	}

	showBookmark := BookmarkShowModel{
		ID:          dbBookmark.ID.String(),
		LastUpdated: dbBookmark.LastUpdated.String(),
		Version:     dbBookmark.Version,
	}

	return ctx.
		Status(fiber.StatusOK).
		JSON(showBookmark)
}
