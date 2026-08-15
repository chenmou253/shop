package app

import (
	"net/http"

	"hbfittings-front/internal/store"

	"github.com/gin-gonic/gin"
)

type App struct {
	store     *store.Store
	uploadDir string
}

func New(cfg Config) (*gin.Engine, error) {
	st, err := store.Open(cfg.DSN)
	if err != nil {
		return nil, err
	}
	a := &App{store: st, uploadDir: cfg.UploadDir}
	r := gin.Default()
	r.Use(cors())
	r.Static("/uploads", cfg.UploadDir)
	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })
	api := r.Group("/api/v1")
	api.Use(noStore())
	api.GET("/bootstrap", a.bootstrap)
	api.GET("/home", a.home)
	api.GET("/categories", a.categories)
	api.GET("/products", a.products)
	api.GET("/products/:slug", a.product)
	api.GET("/articles", a.articles)
	api.GET("/articles/:slug", a.article)
	api.GET("/pages/:slug", a.page)
	api.GET("/videos", a.videos)
	api.GET("/search", a.search)
	api.POST("/inquiries", a.createInquiry)
	api.POST("/subscriptions", a.subscribe)
	return r, nil
}

func noStore() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		c.Next()
	}
}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Accept")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
