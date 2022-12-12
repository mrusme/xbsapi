package bookmarks

import (
	"context"
	"fmt"
	"strings"

	// "github.com/google/uuid"
	"github.com/go-playground/validator/v10"
	"github.com/mrusme/xbsapi/lib"
	"go.uber.org/zap"

	"github.com/gofiber/fiber/v2"
)

const LAST_UPDATED_FORMAT string = "2006-01-02T15:04:05.000Z"

type BookmarkCreateResponse struct {
	Success  bool               `json:"success"`
	Bookmark *BookmarkShowModel `json:"bookmark"`
	Message  string             `json:"message"`
}

// Create godoc
// @Summary      Create a bookmark sync
// @Description  Add a new bookmark sync
// @Tags         bookmarks
// @Accept       json
// @Produce      json
// @Param        bookmark body  BookmarkCreateModel true "Bookmark sync"
// @Success      200  {object}  BookmarkCreateResponse
// @Failure      400  {object}  BookmarkCreateResponse
// @Failure      404  {object}  BookmarkCreateResponse
// @Failure      500  {object}  BookmarkCreateResponse
// @Router       /bookmarks [post]
// @security     BasicAuth
func (h *handler) Create(ctx *fiber.Ctx) error {
	var err error

	if h.config.Service.Status != lib.ServiceStatus(lib.StatusOnline) {
		return ctx.
			Status(fiber.StatusMethodNotAllowed).
			JSON(fiber.Map{
				"success": false,
				"message": "Service not available right now",
			})
	}

	createBookmark := new(BookmarkCreateModel)
	if err = ctx.BodyParser(createBookmark); err != nil {
		h.logger.Debug(
			"Body parsing failed",
			zap.Error(err),
		)
		return ctx.
			Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
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
			JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
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
		ID:          strings.ReplaceAll(dbBookmark.ID.String(), "-", ""),
		LastUpdated: dbBookmark.LastUpdated.Format(LAST_UPDATED_FORMAT),
		Version:     dbBookmark.Version,
	}

	fmt.Println(dbBookmark.ID.String())
	fmt.Printf("%v\n", showBookmark)

	return ctx.
		Status(fiber.StatusOK).
		JSON(showBookmark)
}
