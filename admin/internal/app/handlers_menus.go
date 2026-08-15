package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *App) menusIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "menus/index", a.viewData(c, "菜单管理", "menus"))
}

func (a *App) menusNew(c *gin.Context) {
	data := a.viewData(c, "新增菜单", "menus")
	data["Title"] = "新增菜单"
	data["Mode"] = "create"
	data["ID"] = 0
	c.HTML(http.StatusOK, "menus/form", data)
}

func (a *App) menusEdit(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	data := a.viewData(c, "编辑菜单", "menus")
	data["Title"] = "编辑菜单"
	data["Mode"] = "edit"
	data["ID"] = id
	c.HTML(http.StatusOK, "menus/form", data)
}
