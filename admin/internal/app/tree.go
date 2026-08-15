package app

import (
	"sort"
	"strings"

	"rbac-admin/internal/store"
)

type MenuNode struct {
	Menu     store.Menu `json:"menu"`
	Level    int        `json:"level"`
	Checked  bool       `json:"checked"`
	Children []MenuNode `json:"children"`
}

type PermissionNode struct {
	Permission store.Permission `json:"permission"`
	Level      int              `json:"level"`
	Checked    bool             `json:"checked"`
	Children   []PermissionNode `json:"children"`
}

type FlatMenuNode struct {
	Menu    store.Menu
	Level   int
	Checked bool
}

type FlatPermissionNode struct {
	Permission store.Permission
	Level      int
	Checked    bool
}

func buildMenuTree(menus []store.Menu, selected map[uint]bool) []MenuNode {
	byParent := make(map[uint][]store.Menu)
	for _, menu := range menus {
		byParent[menu.ParentID] = append(byParent[menu.ParentID], menu)
	}
	for parentID := range byParent {
		sort.Slice(byParent[parentID], func(i, j int) bool {
			if byParent[parentID][i].Sort == byParent[parentID][j].Sort {
				return byParent[parentID][i].ID < byParent[parentID][j].ID
			}
			return byParent[parentID][i].Sort < byParent[parentID][j].Sort
		})
	}
	var walk func(parentID uint, level int) []MenuNode
	walk = func(parentID uint, level int) []MenuNode {
		nodes := make([]MenuNode, 0, len(byParent[parentID]))
		for _, menu := range byParent[parentID] {
			nodes = append(nodes, MenuNode{
				Menu:     menu,
				Level:    level,
				Checked:  selected[menu.ID],
				Children: walk(menu.ID, level+1),
			})
		}
		return nodes
	}
	return walk(0, 0)
}

func buildPermissionTree(perms []store.Permission, selected map[uint]bool) []PermissionNode {
	byParent := make(map[uint][]store.Permission)
	for _, perm := range perms {
		byParent[perm.ParentID] = append(byParent[perm.ParentID], perm)
	}
	for parentID := range byParent {
		sort.Slice(byParent[parentID], func(i, j int) bool {
			if byParent[parentID][i].Sort == byParent[parentID][j].Sort {
				return byParent[parentID][i].ID < byParent[parentID][j].ID
			}
			return byParent[parentID][i].Sort < byParent[parentID][j].Sort
		})
	}
	var walk func(parentID uint, level int) []PermissionNode
	walk = func(parentID uint, level int) []PermissionNode {
		nodes := make([]PermissionNode, 0, len(byParent[parentID]))
		for _, perm := range byParent[parentID] {
			nodes = append(nodes, PermissionNode{
				Permission: perm,
				Level:      level,
				Checked:    selected[perm.ID],
				Children:   walk(perm.ID, level+1),
			})
		}
		return nodes
	}
	return walk(0, 0)
}

func flattenMenuTree(nodes []MenuNode) []FlatMenuNode {
	out := make([]FlatMenuNode, 0)
	var walk func([]MenuNode)
	walk = func(items []MenuNode) {
		for _, item := range items {
			out = append(out, FlatMenuNode{Menu: item.Menu, Level: item.Level, Checked: item.Checked})
			walk(item.Children)
		}
	}
	walk(nodes)
	return out
}

func flattenPermissionTree(nodes []PermissionNode) []FlatPermissionNode {
	out := make([]FlatPermissionNode, 0)
	var walk func([]PermissionNode)
	walk = func(items []PermissionNode) {
		for _, item := range items {
			out = append(out, FlatPermissionNode{Permission: item.Permission, Level: item.Level, Checked: item.Checked})
			walk(item.Children)
		}
	}
	walk(nodes)
	return out
}

func selectedMap(ids []uint) map[uint]bool {
	out := make(map[uint]bool, len(ids))
	for _, id := range ids {
		out[id] = true
	}
	return out
}

func parseSelectedCSV(value string) map[uint]bool {
	out := make(map[uint]bool)
	for _, part := range strings.Split(value, ",") {
		id := uintForm(part)
		if id != 0 {
			out[id] = true
		}
	}
	return out
}

func (a *App) allMenus() ([]store.Menu, error) {
	var menus []store.Menu
	err := a.store.DB.Order("sort ASC, id ASC").Find(&menus).Error
	return menus, err
}

func (a *App) allPermissions(menuID uint) ([]store.Permission, error) {
	q := a.store.DB.Preload("Menu").Order("sort ASC, id ASC")
	if menuID != 0 {
		q = q.Where("menu_id = ?", menuID)
	}
	var perms []store.Permission
	err := q.Find(&perms).Error
	return perms, err
}

func (a *App) navMenus(userID uint) ([]MenuNode, error) {
	var menus []store.Menu
	err := a.store.DB.Table("admin_menus m").
		Select("DISTINCT m.*").
		Joins("JOIN admin_role_menus rm ON rm.menu_id = m.id").
		Joins("JOIN admin_user_roles ur ON ur.role_id = rm.role_id").
		Joins("JOIN admin_roles r ON r.id = ur.role_id").
		Where("ur.user_id = ? AND m.visible = ? AND m.status = ? AND r.status = ?", userID, true, true, true).
		Where("m.path NOT IN ?", []string{"/cms/media", "/cms/product-images"}).
		Order("m.sort ASC, m.id ASC").
		Find(&menus).Error
	if err != nil {
		return nil, err
	}
	return buildMenuTree(menus, nil), nil
}
