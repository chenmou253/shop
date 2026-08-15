package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *App) dashboard(c *gin.Context) {
	data := a.viewData(c, "控制台", "dashboard")
	c.HTML(http.StatusOK, "dashboard/index", data)
}
