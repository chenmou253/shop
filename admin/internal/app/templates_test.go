package app

import (
	"bytes"
	"html/template"
	"path/filepath"
	"testing"

	"rbac-admin/internal/store"

	"github.com/gin-gonic/gin"
)

func TestTemplatesParse(t *testing.T) {
	if _, err := parseTestTemplates(); err != nil {
		t.Fatalf("parse templates: %v", err)
	}
}

func TestTemplatesExecute(t *testing.T) {
	tpl, err := parseTestTemplates()
	if err != nil {
		t.Fatalf("parse templates: %v", err)
	}

	menu := store.Menu{ID: 1, Title: "系统管理", Path: "/menus", Visible: true, Status: true, Sort: 10}
	childMenu := store.Menu{ID: 2, ParentID: 1, Title: "用户管理", Path: "/users", Visible: true, Status: true, Sort: 20}
	menuTree := []MenuNode{{Menu: menu, Children: []MenuNode{{Menu: childMenu, Level: 1}}}}
	flatMenus := flattenMenuTree(menuTree)
	role := store.Role{ID: 1, Name: "超级管理员", Code: "super_admin", Status: true}
	user := store.User{ID: 1, Username: "admin", Nickname: "管理员", Email: "admin@example.com", Status: true, Roles: []store.Role{role}}
	perm := store.Permission{ID: 1, MenuID: 2, Name: "用户查看", Code: "user:view", Method: "GET", Path: "/users", Status: true, Menu: childMenu}
	permissionTree := []PermissionNode{{Permission: perm, Checked: true}}
	flatPerms := flattenPermissionTree(permissionTree)

	base := gin.H{
		"Title":       "测试页面",
		"Active":      "dashboard",
		"Username":    "admin",
		"CurrentPath": "/dashboard",
		"NavMenus":    menuTree,
		"Flash":       "",
		"Error":       "",
	}
	pages := map[string]gin.H{
		"auth/login": {
			"Title": "登录",
			"Error": "",
		},
		"error/forbidden": merge(base, gin.H{}),
		"dashboard/index": merge(base, gin.H{
			"UserCount":       int64(1),
			"RoleCount":       int64(1),
			"MenuCount":       int64(2),
			"PermissionCount": int64(1),
		}),
		"users/index": merge(base, gin.H{"Users": []store.User{user}}),
		"users/form": merge(base, gin.H{
			"User":            user,
			"Roles":           []store.Role{role},
			"SelectedRoleIDs": []uint{1},
			"Action":          "/users/1",
			"Mode":            "edit",
		}),
		"roles/index": merge(base, gin.H{"Roles": []store.Role{role}}),
		"roles/form": merge(base, gin.H{
			"Role":                  role,
			"Menus":                 flatMenus,
			"MenuTree":              menuTree,
			"PermissionTree":        permissionTree,
			"SelectedMenuIDs":       []uint{1, 2},
			"SelectedPermissionIDs": []uint{1},
			"Action":                "/roles/1",
			"Mode":                  "edit",
		}),
		"menus/index": merge(base, gin.H{"MenuTree": flatMenus}),
		"menus/form": merge(base, gin.H{
			"Menu":          childMenu,
			"ParentOptions": flatMenus,
			"Action":        "/menus/2",
			"Mode":          "edit",
		}),
		"permissions/index": merge(base, gin.H{"PermissionTree": flatPerms}),
		"permissions/form": merge(base, gin.H{
			"Permission":           perm,
			"Menus":                flatMenus,
			"ParentPermissionTree": permissionTree,
			"Action":               "/permissions/1",
			"Mode":                 "edit",
		}),
		"cms/index": merge(base, gin.H{
			"Resource":    "products",
			"Description": "维护前台产品内容。",
		}),
	}

	for name, data := range pages {
		t.Run(name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := tpl.ExecuteTemplate(&buf, name, data); err != nil {
				t.Fatalf("execute template: %v", err)
			}
			if buf.Len() == 0 {
				t.Fatal("empty output")
			}
		})
	}
}

func parseTestTemplates() (*template.Template, error) {
	pattern := filepath.Join("..", "..", "web", "templates", "*", "*")
	return template.New("").Funcs(templateFuncs()).ParseGlob(pattern)
}

func merge(base gin.H, extra gin.H) gin.H {
	out := gin.H{}
	for key, value := range base {
		out[key] = value
	}
	for key, value := range extra {
		out[key] = value
	}
	return out
}
