package audit

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/ginext"
)

// GET /api/audit ?item_id (optional)
func (h *Handler) GetHistory(c *ginext.Context) {
	if v := c.Query("item_id"); v != "" {
		id, _ := strconv.ParseInt(v, 10, 64)

		items, err := h.service.ListByID(c.Request.Context(), id)
		if err != nil {
			h.log.Error().Err(err).Str("item_id", v).Msg("failed to get item history")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, items)
		return
	}
	items, err := h.service.ListAll(c.Request.Context())
	if err != nil {
		h.log.Error().Err(err).Msg("failed to get item history")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GET /api/audit/export
func (h *Handler) ExportHistoryCSV(c *ginext.Context) {
	// minimal CSV export
	out, err := h.service.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=history.csv")
	c.Writer.WriteString("id,item_id,action,user_id,created_at,old_value,new_value\n")
	for _, r := range out {
		c.Writer.WriteString(fmt.Sprintf("%d,%d,%s,%d,%s,%q,%q\n", r.ID, r.ItemID, r.Action, r.UserID, r.CreatedAt.Format(time.RFC3339), r.OldValue, r.NewValue))
	}
}
