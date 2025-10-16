package items

import (
	"errors"
	"github.com/K1la/warehouse-control/internal/dto"
	itemRepo "github.com/K1la/warehouse-control/internal/repository/item"
	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/ginext"
	"net/http"
)

// POST /item
func (h *Handler) Create(c *ginext.Context) {
	var req dto.CreateItem
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error().Err(err).Msg("failed to bind create item request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.log.Info().Interface("item", req).Msg("parsed create item")

	item, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Msg("failed to create item from service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.log.Info().Interface("item", item).Msg("created item")
	c.JSON(http.StatusOK, item)
}

// GET /item
func (h *Handler) GetAll(c *ginext.Context) {
	var params dto.GetItemsParams
	if err := c.ShouldBindQuery(&params); err != nil {
		h.log.Error().Err(err).Msg("failed to bind get item params")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	items, err := h.service.GetAll(c.Request.Context(), params)
	if err != nil {
		h.log.Error().Err(err).Msg("failed to get item from service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GET /item/:id
func (h *Handler) GetByID(c *ginext.Context) {
	id := c.Param("id")

	item, err := h.service.GetByID(c.Request.Context(), id)
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

// PUT /item/:id
func (h *Handler) Update(c *ginext.Context) {
	id := c.Param("id")

	var req dto.UpdateItem
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error().Err(err).Msg("failed to bind update item request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, itemRepo.ErrItemNotFound) {
			h.log.Error().Err(err).Msg("item not found")
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.log.Error().Err(err).Msg("failed to update item from service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *Handler) Delete(c *ginext.Context) {
	id := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		h.log.Error().Err(err).Msg("failed to delete item from service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
