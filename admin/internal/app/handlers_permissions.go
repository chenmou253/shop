package app

import (
	"net/http"

	"rbac-admin/internal/store"

	"github.com/gin-gonic/gin"
)

func (a *App) permissionsIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "permissions/index", a.viewData(c, "权限管理", "permissions"))
}

func (a *App) permissionsNew(c *gin.Context) {
	data := a.viewData(c, "新增权限", "permissions")
	data["Title"] = "新增权限"
	data["Mode"] = "create"
	data["ID"] = 0
	c.HTML(http.StatusOK, "permissions/form", data)
}

func (a *App) permissionsEdit(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	data := a.viewData(c, "编辑权限", "permissions")
	data["Title"] = "编辑权限"
	data["Mode"] = "edit"
	data["ID"] = id
	c.HTML(http.StatusOK, "permissions/form", data)
}

func (a *App) parentPermissions(menuID, currentID uint) ([]store.Permission, error) {
	q := a.store.DB.Preload("Menu").Order("sort ASC, id ASC")
	if menuID != 0 {
		q = q.Where("menu_id = ?", menuID)
	}
	if currentID != 0 {
		q = q.Where("id <> ?", currentID)
	}
	var perms []store.Permission
	err := q.Find(&perms).Error
	return perms, err
}
