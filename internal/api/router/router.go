package router

import (
	"github.com/K1la/warehouse-control/internal/middleware"
	"net/http"
	"path/filepath"
	"strings"

	auditH "github.com/K1la/warehouse-control/internal/api/handlers/audit"
	authH "github.com/K1la/warehouse-control/internal/api/handlers/auth"
	itemH "github.com/K1la/warehouse-control/internal/api/handlers/items"

	"github.com/wb-go/wbf/ginext"
)

func New(itH *itemH.Handler, atH *authH.Handler, auditH *auditH.Handler, authMw *middleware.AuthMW) *ginext.Engine {
	e := ginext.New("")
	e.Use(ginext.Recovery(), ginext.Logger())

	// API routes
	api := e.Group("/api")
	{
		// Auth routes
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/login", atH.Login)
			authGroup.POST("/register", atH.Register)
		}

		// Audit routes
		auditGroup := api.Group("/audit")
		auditGroup.Use(authMw.RequireAuth())
		{
			auditGroup.GET("", auditH.GetHistory)
			auditGroup.GET("/export", authMw.RequireRole("admin", "manager"), auditH.ExportHistoryCSV)
		}

		// Items routes
		itemsGroup := api.Group("/items")
		{
			itemsGroup.GET("", itH.GetAll)
			itemsGroup.GET("/:id", itH.GetByID)
			itemsGroup.POST("", authMw.RequireRole("admin", "manager"), itH.Create)
			itemsGroup.PUT("/:id", authMw.RequireRole("admin", "manager"), itH.Update)
			itemsGroup.DELETE("/:id", authMw.RequireRole("admin"), itH.Delete)
		}
	}

	// Frontend: serve files from ./web without conflicting wildcard
	e.NoRoute(func(c *ginext.Context) {
		if c.Request.URL.Path == "/" {
			http.ServeFile(c.Writer, c.Request, "./web/index.html")
			return
		}
		// Serve only files under /web/ directly from disk
		if strings.HasPrefix(c.Request.URL.Path, "/web/") {
			safe := filepath.Clean("." + c.Request.URL.Path)
			http.ServeFile(c.Writer, c.Request, safe)
			return
		}
		c.Status(http.StatusNotFound)
	})

	return e
}
