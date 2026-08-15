package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *App) usersIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "users/index", a.viewData(c, "用户管理", "users"))
}

func (a *App) usersNew(c *gin.Context) {
	data := a.viewData(c, "新增用户", "users")
	data["Title"] = "新增用户"
	data["Mode"] = "create"
	data["ID"] = 0
	c.HTML(http.StatusOK, "users/form", data)
}

func (a *App) usersEdit(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	data := a.viewData(c, "编辑用户", "users")
	data["Title"] = "编辑用户"
	data["Mode"] = "edit"
	data["ID"] = id
	c.HTML(http.StatusOK, "users/form", data)
}
