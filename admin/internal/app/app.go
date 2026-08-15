package app

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"rbac-admin/internal/session"
	"rbac-admin/internal/store"

	"github.com/gin-gonic/gin"
)

type App struct {
	cfg      Config
	store    *store.Store
	sessions *session.Manager
}

func New(cfg Config) (*gin.Engine, error) {
	st, err := store.Open(cfg.DSN)
	if err != nil {
		return nil, err
	}

	a := &App{
		cfg:      cfg,
		store:    st,
		sessions: session.NewManager(cfg.SessionSecret),
	}

	r := gin.Default()
	r.SetFuncMap(templateFuncs())
	r.LoadHTMLGlob("web/templates/**/*")
	r.Use(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/static/") {
			c.Header("Cache-Control", "no-store, max-age=0")
		}
		c.Next()
	})
	r.Static("/static", "web/static")
	r.Static("/uploads", cfg.UploadDir)
	a.routes(r)
	return r, nil
}

func templateFuncs() template.FuncMap {
	return template.FuncMap{
		"joinUint":      joinUint,
		"containsUint":  containsUint,
		"indent":        indent,
		"methodBadge":   methodBadge,
		"statusText":    statusText,
		"statusBadge":   statusBadge,
		"activeSection": activeSection,
		"roleNames":     roleNames,
	}
}

func (a *App) routes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) { c.Redirect(http.StatusFound, "/dashboard") })
	r.GET("/login", a.loginPage)

	api := r.Group("/api")
	api.POST("/login", a.apiLogin)

	authedAPI := api.Group("/")
	authedAPI.Use(a.authRequiredAPI())
	authedAPI.GET("/me", a.apiMe)
	authedAPI.POST("/logout", a.apiLogout)
	authedAPI.GET("/dashboard", a.apiDashboard)
	authedAPI.GET("/nav", a.apiNav)

	authedAPI.GET("/users", a.permissionRequiredAPI("user:view"), a.apiUsersIndex)
	authedAPI.POST("/users", a.permissionRequiredAPI("user:create"), a.apiUsersCreate)
	authedAPI.GET("/users/:id", a.permissionRequiredAPI("user:update"), a.apiUsersShow)
	authedAPI.PUT("/users/:id", a.permissionRequiredAPI("user:update"), a.apiUsersUpdate)
	authedAPI.DELETE("/users/:id", a.permissionRequiredAPI("user:delete"), a.apiUsersDelete)

	authedAPI.GET("/roles", a.permissionRequiredAPI("role:view"), a.apiRolesIndex)
	authedAPI.POST("/roles", a.permissionRequiredAPI("role:create"), a.apiRolesCreate)
	authedAPI.GET("/roles/:id", a.permissionRequiredAPI("role:update"), a.apiRolesShow)
	authedAPI.PUT("/roles/:id", a.permissionRequiredAPI("role:update"), a.apiRolesUpdate)
	authedAPI.DELETE("/roles/:id", a.permissionRequiredAPI("role:delete"), a.apiRolesDelete)

	authedAPI.GET("/menus", a.permissionRequiredAPI("menu:view"), a.apiMenusIndex)
	authedAPI.POST("/menus", a.permissionRequiredAPI("menu:create"), a.apiMenusCreate)
	authedAPI.GET("/menus/:id", a.permissionRequiredAPI("menu:update"), a.apiMenusShow)
	authedAPI.PUT("/menus/:id", a.permissionRequiredAPI("menu:update"), a.apiMenusUpdate)
	authedAPI.DELETE("/menus/:id", a.permissionRequiredAPI("menu:delete"), a.apiMenusDelete)

	authedAPI.GET("/permissions", a.permissionRequiredAPI("permission:view"), a.apiPermissionsIndex)
	authedAPI.POST("/permissions", a.permissionRequiredAPI("permission:create"), a.apiPermissionsCreate)
	authedAPI.GET("/permissions/:id", a.permissionRequiredAPI("permission:update"), a.apiPermissionsShow)
	authedAPI.PUT("/permissions/:id", a.permissionRequiredAPI("permission:update"), a.apiPermissionsUpdate)
	authedAPI.DELETE("/permissions/:id", a.permissionRequiredAPI("permission:delete"), a.apiPermissionsDelete)

	authedAPI.GET("/menus/tree", a.apiMenusTree)
	authedAPI.GET("/permissions/tree", a.apiPermissionsTree)
	authedAPI.GET("/role-access", a.apiRoleAccess)
	authedAPI.GET("/permission-parents", a.apiPermissionParents)

	authedAPI.GET("/cms/:resource/meta", a.cmsPermission("view"), a.apiCMSMeta)
	authedAPI.GET("/cms/:resource", a.cmsPermission("view"), a.apiCMSIndex)
	authedAPI.POST("/cms/:resource", a.cmsPermission("create"), a.apiCMSCreate)
	authedAPI.POST("/cms/:resource/upload", a.cmsUploadPermission(), a.apiImageUpload)
	authedAPI.GET("/cms/:resource/:id", a.cmsPermission("update"), a.apiCMSShow)
	authedAPI.PUT("/cms/:resource/:id", a.cmsPermission("update"), a.apiCMSUpdate)
	authedAPI.DELETE("/cms/:resource/:id", a.cmsPermission("delete"), a.apiCMSDelete)

	admin := r.Group("/")
	admin.Use(a.authRequired())
	admin.GET("/dashboard", a.dashboard)

	admin.GET("/users", a.permissionRequired("user:view"), a.usersIndex)
	admin.GET("/users/new", a.permissionRequired("user:create"), a.usersNew)
	admin.GET("/users/:id/edit", a.permissionRequired("user:update"), a.usersEdit)

	admin.GET("/roles", a.permissionRequired("role:view"), a.rolesIndex)
	admin.GET("/roles/new", a.permissionRequired("role:create"), a.rolesNew)
	admin.GET("/roles/:id/edit", a.permissionRequired("role:update"), a.rolesEdit)

	admin.GET("/menus", a.permissionRequired("menu:view"), a.menusIndex)
	admin.GET("/menus/new", a.permissionRequired("menu:create"), a.menusNew)
	admin.GET("/menus/:id/edit", a.permissionRequired("menu:update"), a.menusEdit)

	admin.GET("/permissions", a.permissionRequired("permission:view"), a.permissionsIndex)
	admin.GET("/permissions/new", a.permissionRequired("permission:create"), a.permissionsNew)
	admin.GET("/permissions/:id/edit", a.permissionRequired("permission:update"), a.permissionsEdit)

	admin.GET("/cms/:resource", a.cmsPage)
}

func (a *App) authRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		sess, ok := a.sessions.Get(c)
		if !ok {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		c.Set("user_id", sess.UserID)
		c.Set("username", sess.Username)
		c.Next()
	}
}

func (a *App) permissionRequired(code string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := currentUserID(c)
		if !a.store.UserHasPermission(userID, code) {
			c.HTML(http.StatusForbidden, "error/forbidden", a.viewData(c, "无权限", ""))
			c.Abort()
			return
		}
		c.Next()
	}
}

func (a *App) authRequiredAPI() gin.HandlerFunc {
	return func(c *gin.Context) {
		sess, ok := a.sessions.Get(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
			c.Abort()
			return
		}
		c.Set("user_id", sess.UserID)
		c.Set("username", sess.Username)
		c.Next()
	}
}

func (a *App) permissionRequiredAPI(code string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := currentUserID(c)
		if !a.store.UserHasPermission(userID, code) {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权限"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func (a *App) viewData(c *gin.Context, title, active string) gin.H {
	username, _ := c.Get("username")
	return gin.H{
		"Title":       title,
		"Active":      active,
		"Username":    username,
		"CurrentPath": c.Request.URL.Path,
		"NavMenus":    []MenuNode{},
		"Flash":       c.Query("flash"),
		"Error":       c.Query("error"),
	}
}

func currentUserID(c *gin.Context) uint {
	value, ok := c.Get("user_id")
	if !ok {
		return 0
	}
	id, ok := value.(uint)
	if !ok {
		return 0
	}
	return id
}

func parseIDParam(c *gin.Context, name string) (uint, bool) {
	value, err := strconv.ParseUint(c.Param(name), 10, 64)
	if err != nil || value == 0 {
		c.String(http.StatusBadRequest, "invalid id")
		return 0, false
	}
	return uint(value), true
}

func joinUint(ids []uint) string {
	parts := make([]string, 0, len(ids))
	for _, id := range ids {
		parts = append(parts, strconv.FormatUint(uint64(id), 10))
	}
	return strings.Join(parts, ",")
}

func containsUint(ids []uint, id uint) bool {
	for _, item := range ids {
		if item == id {
			return true
		}
	}
	return false
}

func indent(level int) template.HTML {
	if level <= 0 {
		return ""
	}
	return template.HTML(strings.Repeat("&nbsp;&nbsp;&nbsp;&nbsp;", level))
}

func methodBadge(method string) string {
	if method == "" {
		return "secondary"
	}
	switch strings.ToUpper(method) {
	case http.MethodGet:
		return "success"
	case http.MethodPost:
		return "primary"
	case http.MethodPut, http.MethodPatch:
		return "warning"
	case http.MethodDelete:
		return "danger"
	default:
		return "secondary"
	}
}

func statusText(status bool) string {
	if status {
		return "启用"
	}
	return "禁用"
}

func statusBadge(status bool) string {
	if status {
		return "success"
	}
	return "secondary"
}

func activeSection(active, section string) bool {
	return strings.HasPrefix(active, section)
}

func atoiDefault(value string, fallback int) int {
	n, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return fallback
	}
	return n
}

func uintForm(value string) uint {
	n, err := strconv.ParseUint(strings.TrimSpace(value), 10, 64)
	if err != nil {
		return 0
	}
	return uint(n)
}

func boolForm(value string) bool {
	return value == "1" || value == "on" || value == "true"
}

func fail(c *gin.Context, templateName string, data gin.H, err error) {
	data["Error"] = fmt.Sprintf("操作失败：%v", err)
	c.HTML(http.StatusOK, templateName, data)
}
