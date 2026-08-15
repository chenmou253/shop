package app

import (
	"net/http"

	"rbac-admin/internal/store"

	"github.com/gin-gonic/gin"
)

func (a *App) apiMenusTree(c *gin.Context) {
	menus, err := a.allMenus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"tree": buildMenuTree(menus, parseSelectedCSV(c.Query("selected"))),
	})
}

func (a *App) apiPermissionsTree(c *gin.Context) {
	menuID := uintForm(c.Query("menu_id"))
	perms, err := a.allPermissions(menuID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"tree": buildPermissionTree(perms, parseSelectedCSV(c.Query("selected"))),
	})
}

func (a *App) apiRoleAccess(c *gin.Context) {
	roleID := uintForm(c.Query("role_id"))
	menuID := uintForm(c.Query("menu_id"))
	var selectedMenus, selectedPerms []uint
	if roleID != 0 {
		selectedMenus, _ = a.store.MenuIDsForRole(roleID)
		selectedPerms, _ = a.store.PermissionIDsForRole(roleID)
	}
	if values, ok := c.GetQueryArray("selected_menus"); ok {
		selectedMenus = store.ParseUintList(values)
	}
	if values, ok := c.GetQueryArray("selected_permissions"); ok {
		selectedPerms = store.ParseUintList(values)
	}

	menus, err := a.allMenus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	perms, err := a.allPermissions(menuID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"menus":                   buildMenuTree(menus, selectedMap(selectedMenus)),
		"permissions":             buildPermissionTree(perms, selectedMap(selectedPerms)),
		"selected_menu_ids":       selectedMenus,
		"selected_permission_ids": selectedPerms,
	})
}

func (a *App) apiPermissionParents(c *gin.Context) {
	menuID := uintForm(c.Query("menu_id"))
	currentID := uintForm(c.Query("current_id"))
	checkedID := uintForm(c.Query("checked"))
	perms, err := a.parentPermissions(menuID, currentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"tree": buildPermissionTree(perms, selectedMap([]uint{checkedID})),
	})
}

func roleNames(roles []store.Role) string {
	if len(roles) == 0 {
		return "-"
	}
	names := make([]string, 0, len(roles))
	for _, role := range roles {
		names = append(names, role.Name)
	}
	return joinStrings(names, "、")
}

func joinStrings(values []string, sep string) string {
	if len(values) == 0 {
		return ""
	}
	out := values[0]
	for i := 1; i < len(values); i++ {
		out += sep + values[i]
	}
	return out
}
