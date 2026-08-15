package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *App) rolesIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "roles/index", a.viewData(c, "角色管理", "roles"))
}

func (a *App) rolesNew(c *gin.Context) {
	data := a.viewData(c, "新增角色", "roles")
	data["Title"] = "新增角色"
	data["Mode"] = "create"
	data["ID"] = 0
	c.HTML(http.StatusOK, "roles/form", data)
}

func (a *App) rolesEdit(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	data := a.viewData(c, "编辑角色", "roles")
	data["Title"] = "编辑角色"
	data["Mode"] = "edit"
	data["ID"] = id
	c.HTML(http.StatusOK, "roles/form", data)
}
