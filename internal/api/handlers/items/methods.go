package items

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/K1la/warehouse-control/internal/dto"
	itemRepo "github.com/K1la/warehouse-control/internal/repository/item"
	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/ginext"
)

// POST /api/items
func (h *Handler) Create(c *ginext.Context) {
	var req dto.CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error().Err(err).Msg("failed to bind create item request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get user_id from ctx (set vy auth middleware)
	uid := c.GetInt64("user_id")

	h.log.Info().Interface("req", req).Int64("user_id", uid).
		Msg("parsed created item and user id")

	id, err := h.service.Create(c.Request.Context(), uid, req)
	if err != nil {
		h.log.Error().Err(err).Msg("failed to create item from service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.log.Info().Interface("item id", id).Msg("created item")
	c.JSON(http.StatusCreated, ginext.H{"id": id})
}

// GET /api/items
func (h *Handler) GetAll(c *ginext.Context) {
	items, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		h.log.Error().Err(err).Msg("failed to get item from service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GET /api/items/:id
func (h *Handler) GetByID(c *ginext.Context) {
	itemIDStr := c.Param("id")
	itemID, err := strconv.ParseInt(itemIDStr, 10, 64)
	if err != nil {
		h.log.Error().Err(err).Msg("failed to parse item id")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.service.GetByID(c.Request.Context(), itemID)
	if err != nil {
		if errors.Is(err, itemRepo.ErrItemNotFound) {
			h.log.Error().Err(err).Msg("item not found")
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.log.Error().Err(err).Msg("failed to get item from service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

// PUT /api/items/:id
func (h *Handler) Update(c *ginext.Context) {
	var req dto.UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error().Err(err).Msg("failed to bind update item request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, itemID, ok := h.getUserAndItemIDFromContext(c)
	if !ok {
		return
	}

	err := h.service.Update(c.Request.Context(), userID, itemID, req)
	if err != nil {
		if errors.Is(err, itemRepo.ErrItemNotFound) {
			h.log.Error().Err(err).Msg("item not found")
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.log.Error().Err(err).Msg("failed to get item from service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.log.Info().Int64("id", itemID).Interface("item", req).
		Msg("updated item")
	c.JSON(http.StatusOK, ginext.H{"id": itemID})
}

// DELETE /api/items/:id
func (h *Handler) Delete(c *ginext.Context) {
	userID, itemID, ok := h.getUserAndItemIDFromContext(c)
	if !ok {
		return
	}

	if err := h.service.Delete(c.Request.Context(), userID, itemID); err != nil {
		h.log.Error().Err(err).Msg("failed to delete item from service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ginext.H{"id": itemID})
}

func (h *Handler) getUserAndItemIDFromContext(c *ginext.Context) (int64, int64, bool) {
	itemIDStr := c.Param("id")
	itemID, err := strconv.ParseInt(itemIDStr, 10, 64)
	if err != nil {
		h.log.Error().Err(err).Msg("failed to parse item id")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return 0, 0, false
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		h.log.Error().Interface("userIDVal", userIDVal).
			Msg("failed to get user id from ctx")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "userID not found in context"})
		return 0, 0, false
	}

	userID, ok := userIDVal.(int64)
	if !ok {
		h.log.Error().Interface("userid", userID).
			Msgf("failed to get user id from ctx, ok=%v", ok)
		c.JSON(http.StatusUnauthorized, fmt.Errorf("invalid userID type"))
		return 0, 0, false
	}

	return userID, itemID, true
}
