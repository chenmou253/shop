package app

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"rbac-admin/internal/store"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type userDTO struct {
	ID       uint      `json:"id"`
	Username string    `json:"username"`
	Nickname string    `json:"nickname"`
	Email    string    `json:"email"`
	Status   bool      `json:"status"`
	Roles    []roleDTO `json:"roles"`
	RoleIDs  []uint    `json:"role_ids,omitempty"`
}

type roleDTO struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	Remark   string `json:"remark"`
	Status   bool   `json:"status"`
	MenuIDs  []uint `json:"menu_ids,omitempty"`
	PermIDs  []uint `json:"permission_ids,omitempty"`
	UserText string `json:"user_text,omitempty"`
}

type menuDTO struct {
	ID       uint   `json:"id"`
	ParentID uint   `json:"parent_id"`
	Title    string `json:"title"`
	Icon     string `json:"icon"`
	Path     string `json:"path"`
	Sort     int    `json:"sort"`
	Visible  bool   `json:"visible"`
	Status   bool   `json:"status"`
}

type permissionDTO struct {
	ID       uint    `json:"id"`
	ParentID uint    `json:"parent_id"`
	MenuID   uint    `json:"menu_id"`
	Name     string  `json:"name"`
	Code     string  `json:"code"`
	Method   string  `json:"method"`
	Path     string  `json:"path"`
	Sort     int     `json:"sort"`
	Status   bool    `json:"status"`
	Menu     menuDTO `json:"menu"`
}

type userRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Status   bool   `json:"status"`
	RoleIDs  []uint `json:"role_ids"`
}

type roleRequest struct {
	Name          string `json:"name"`
	Code          string `json:"code"`
	Remark        string `json:"remark"`
	Status        bool   `json:"status"`
	MenuIDs       []uint `json:"menu_ids"`
	PermissionIDs []uint `json:"permission_ids"`
}

type menuRequest struct {
	ParentID uint   `json:"parent_id"`
	Title    string `json:"title"`
	Icon     string `json:"icon"`
	Path     string `json:"path"`
	Sort     int    `json:"sort"`
	Visible  bool   `json:"visible"`
	Status   bool   `json:"status"`
}

type permissionRequest struct {
	ParentID uint   `json:"parent_id"`
	MenuID   uint   `json:"menu_id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	Method   string `json:"method"`
	Path     string `json:"path"`
	Sort     int    `json:"sort"`
	Status   bool   `json:"status"`
}

func (a *App) apiDashboard(c *gin.Context) {
	var userCount, roleCount, menuCount, permissionCount int64
	var productCount, articleCount, newInquiryCount int64
	a.store.DB.Table("admin_users").Count(&userCount)
	a.store.DB.Table("admin_roles").Count(&roleCount)
	a.store.DB.Table("admin_menus").Count(&menuCount)
	a.store.DB.Table("admin_permissions").Count(&permissionCount)
	a.store.DB.Table("products").Where("status = ?", true).Count(&productCount)
	a.store.DB.Table("content_articles").Where("status = ?", true).Count(&articleCount)
	a.store.DB.Table("inquiries").Where("lead_status = ?", "new").Count(&newInquiryCount)
	c.JSON(http.StatusOK, gin.H{
		"user_count":       userCount,
		"role_count":       roleCount,
		"menu_count":       menuCount,
		"permission_count": permissionCount,
		"product_count":    productCount,
		"article_count":    articleCount,
		"inquiry_count":    newInquiryCount,
	})
}

func (a *App) apiUsersIndex(c *gin.Context) {
	var users []store.User
	if err := a.store.DB.Preload("Roles").Order("id DESC").Find(&users).Error; err != nil {
		jsonError(c, http.StatusInternalServerError, err)
		return
	}
	out := make([]userDTO, 0, len(users))
	for _, user := range users {
		out = append(out, makeUserDTO(user, nil))
	}
	c.JSON(http.StatusOK, gin.H{"users": out})
}

func (a *App) apiUsersShow(c *gin.Context) {
	id, ok := parseAPIIDParam(c, "id")
	if !ok {
		return
	}
	var user store.User
	if err := a.store.DB.Preload("Roles").First(&user, id).Error; err != nil {
		jsonError(c, http.StatusNotFound, err)
		return
	}
	roleIDs, _ := a.store.RoleIDsForUser(user.ID)
	c.JSON(http.StatusOK, gin.H{"user": makeUserDTO(user, roleIDs)})
}

func (a *App) apiUsersCreate(c *gin.Context) {
	var req userRequest
	if !bindJSON(c, &req) {
		return
	}
	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)
	if req.Username == "" || req.Password == "" {
		jsonMessage(c, http.StatusBadRequest, "用户名和密码不能为空")
		return
	}
	salt, err := store.RandomSalt()
	if err != nil {
		jsonError(c, http.StatusInternalServerError, err)
		return
	}
	user := store.User{
		Username:     req.Username,
		PasswordSalt: salt,
		PasswordHash: store.HashPassword(salt, req.Password),
		Nickname:     strings.TrimSpace(req.Nickname),
		Email:        strings.TrimSpace(req.Email),
		Status:       req.Status,
	}
	if err := a.store.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		return a.store.ReplaceUserRoles(tx, user.ID, req.RoleIDs)
	}); err != nil {
		jsonError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "用户已创建", "user": makeUserDTO(user, req.RoleIDs)})
}

func (a *App) apiUsersUpdate(c *gin.Context) {
	id, ok := parseAPIIDParam(c, "id")
	if !ok {
		return
	}
	var req userRequest
	if !bindJSON(c, &req) {
		return
	}
	var user store.User
	if err := a.store.DB.First(&user, id).Error; err != nil {
		jsonError(c, http.StatusNotFound, err)
		return
	}
	user.Username = strings.TrimSpace(req.Username)
	user.Nickname = strings.TrimSpace(req.Nickname)
	user.Email = strings.TrimSpace(req.Email)
	user.Status = req.Status
	if user.Username == "" {
		jsonMessage(c, http.StatusBadRequest, "用户名不能为空")
		return
	}
	if password := strings.TrimSpace(req.Password); password != "" {
		salt, err := store.RandomSalt()
		if err != nil {
			jsonError(c, http.StatusInternalServerError, err)
			return
		}
		user.PasswordSalt = salt
		user.PasswordHash = store.HashPassword(salt, password)
	}
	if err := a.store.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&user).Error; err != nil {
			return err
		}
		return a.store.ReplaceUserRoles(tx, user.ID, req.RoleIDs)
	}); err != nil {
		jsonError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "用户已更新", "user": makeUserDTO(user, req.RoleIDs)})
}

func (a *App) apiUsersDelete(c *gin.Context) {
	id, ok := parseAPIIDParam(c, "id")
	if !ok {
		return
	}
	err := a.store.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", id).Delete(&store.UserRole{}).Error; err != nil {
			return err
		}
		return tx.Delete(&store.User{}, id).Error
	})
	if err != nil {
		jsonError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "用户已删除"})
}

func (a *App) apiRolesIndex(c *gin.Context) {
	var roles []store.Role
	if err := a.store.DB.Order("id DESC").Find(&roles).Error; err != nil {
		jsonError(c, http.StatusInternalServerError, err)
		return
	}
	out := make([]roleDTO, 0, len(roles))
	for _, role := range roles {
		out = append(out, makeRoleDTO(role, nil, nil))
	}
	c.JSON(http.StatusOK, gin.H{"roles": out})
}

func (a *App) apiRolesShow(c *gin.Context) {
	id, ok := parseAPIIDParam(c, "id")
	if !ok {
		return
	}
	var role store.Role
	if err := a.store.DB.First(&role, id).Error; err != nil {
		jsonError(c, http.StatusNotFound, err)
		return
	}
	menuIDs, _ := a.store.MenuIDsForRole(role.ID)
	permIDs, _ := a.store.PermissionIDsForRole(role.ID)
	c.JSON(http.StatusOK, gin.H{"role": makeRoleDTO(role, menuIDs, permIDs)})
}

func (a *App) apiRolesCreate(c *gin.Context) {
	var req roleRequest
	if !bindJSON(c, &req) {
		return
	}
	role := store.Role{
		Name:   strings.TrimSpace(req.Name),
		Code:   strings.TrimSpace(req.Code),
		Remark: strings.TrimSpace(req.Remark),
		Status: req.Status,
	}
	if role.Name == "" || role.Code == "" {
		jsonMessage(c, http.StatusBadRequest, "角色名称和标识不能为空")
		return
	}
	if err := a.store.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&role).Error; err != nil {
			return err
		}
		return a.store.ReplaceRoleAccess(tx, role.ID, req.MenuIDs, req.PermissionIDs)
	}); err != nil {
		jsonError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "角色已创建", "role": makeRoleDTO(role, req.MenuIDs, req.PermissionIDs)})
}

func (a *App) apiRolesUpdate(c *gin.Context) {
	id, ok := parseAPIIDParam(c, "id")
	if !ok {
		return
	}
	var req roleRequest
	if !bindJSON(c, &req) {
		return
	}
	var role store.Role
	if err := a.store.DB.First(&role, id).Error; err != nil {
		jsonError(c, http.StatusNotFound, err)
		return
	}
	role.Name = strings.TrimSpace(req.Name)
	role.Code = strings.TrimSpace(req.Code)
	role.Remark = strings.TrimSpace(req.Remark)
	role.Status = req.Status
	if role.Name == "" || role.Code == "" {
		jsonMessage(c, http.StatusBadRequest, "角色名称和标识不能为空")
		return
	}
	if err := a.store.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&role).Error; err != nil {
			return err
		}
		return a.store.ReplaceRoleAccess(tx, role.ID, req.MenuIDs, req.PermissionIDs)
	}); err != nil {
		jsonError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "角色已更新", "role": makeRoleDTO(role, req.MenuIDs, req.PermissionIDs)})
}

func (a *App) apiRolesDelete(c *gin.Context) {
	id, ok := parseAPIIDParam(c, "id")
	if !ok {
		return
	}
	err := a.store.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", id).Delete(&store.UserRole{}).Error; err != nil {
			return err
		}
		if err := tx.Where("role_id = ?", id).Delete(&store.RoleMenu{}).Error; err != nil {
			return err
		}
		if err := tx.Where("role_id = ?", id).Delete(&store.RolePermission{}).Error; err != nil {
			return err
		}
		return tx.Delete(&store.Role{}, id).Error
	})
	if err != nil {
		jsonError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "角色已删除"})
}

func (a *App) apiMenusIndex(c *gin.Context) {
	menus, err := a.allMenus()
	if err != nil {
		jsonError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"menus": makeMenuDTOs(menus),
		"tree":  buildMenuTree(menus, nil),
	})
}

func (a *App) apiMenusShow(c *gin.Context) {
	id, ok := parseAPIIDParam(c, "id")
	if !ok {
		return
	}
	var menu store.Menu
	if err := a.store.DB.First(&menu, id).Error; err != nil {
		jsonError(c, http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"menu": makeMenuDTO(menu)})
}

func (a *App) apiMenusCreate(c *gin.Context) {
	var req menuRequest
	if !bindJSON(c, &req) {
		return
	}
	menu := store.Menu{
		ParentID: req.ParentID,
		Title:    strings.TrimSpace(req.Title),
		Icon:     strings.TrimSpace(req.Icon),
		Path:     strings.TrimSpace(req.Path),
		Sort:     req.Sort,
		Visible:  req.Visible,
		Status:   req.Status,
	}
	if menu.Title == "" {
		jsonMessage(c, http.StatusBadRequest, "菜单名称不能为空")
		return
	}
	if err := a.store.DB.Create(&menu).Error; err != nil {
		jsonError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "菜单已创建", "menu": makeMenuDTO(menu)})
}

func (a *App) apiMenusUpdate(c *gin.Context) {
	id, ok := parseAPIIDParam(c, "id")
	if !ok {
		return
	}
	var req menuRequest
	if !bindJSON(c, &req) {
		return
	}
	var menu store.Menu
	if err := a.store.DB.First(&menu, id).Error; err != nil {
		jsonError(c, http.StatusNotFound, err)
		return
	}
	if req.ParentID == id {
		jsonMessage(c, http.StatusBadRequest, "上级菜单不能选择自身")
		return
	}
	menu.ParentID = req.ParentID
	menu.Title = strings.TrimSpace(req.Title)
	menu.Icon = strings.TrimSpace(req.Icon)
	menu.Path = strings.TrimSpace(req.Path)
	menu.Sort = req.Sort
	menu.Visible = req.Visible
	menu.Status = req.Status
	if menu.Title == "" {
		jsonMessage(c, http.StatusBadRequest, "菜单名称不能为空")
		return
	}
	if err := a.store.DB.Save(&menu).Error; err != nil {
		jsonError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "菜单已更新", "menu": makeMenuDTO(menu)})
}

func (a *App) apiMenusDelete(c *gin.Context) {
	id, ok := parseAPIIDParam(c, "id")
	if !ok {
		return
	}
	err := a.store.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&store.Menu{}).Where("parent_id = ?", id).Update("parent_id", 0).Error; err != nil {
			return err
		}
		if err := tx.Model(&store.Permission{}).Where("menu_id = ?", id).Update("menu_id", 0).Error; err != nil {
			return err
		}
		if err := tx.Where("menu_id = ?", id).Delete(&store.RoleMenu{}).Error; err != nil {
			return err
		}
		return tx.Delete(&store.Menu{}, id).Error
	})
	if err != nil {
		jsonError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "菜单已删除"})
}

func (a *App) apiPermissionsIndex(c *gin.Context) {
	perms, err := a.allPermissions(0)
	if err != nil {
		jsonError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"permissions": makePermissionDTOs(perms),
		"tree":        buildPermissionTree(perms, nil),
	})
}

func (a *App) apiPermissionsShow(c *gin.Context) {
	id, ok := parseAPIIDParam(c, "id")
	if !ok {
		return
	}
	var perm store.Permission
	if err := a.store.DB.Preload("Menu").First(&perm, id).Error; err != nil {
		jsonError(c, http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"permission": makePermissionDTO(perm)})
}

func (a *App) apiPermissionsCreate(c *gin.Context) {
	var req permissionRequest
	if !bindJSON(c, &req) {
		return
	}
	perm := store.Permission{
		ParentID: req.ParentID,
		MenuID:   req.MenuID,
		Name:     strings.TrimSpace(req.Name),
		Code:     strings.TrimSpace(req.Code),
		Method:   strings.ToUpper(strings.TrimSpace(req.Method)),
		Path:     strings.TrimSpace(req.Path),
		Sort:     req.Sort,
		Status:   req.Status,
	}
	if perm.Method == "" {
		perm.Method = http.MethodGet
	}
	if perm.Name == "" || perm.Code == "" {
		jsonMessage(c, http.StatusBadRequest, "权限名称和权限标识不能为空")
		return
	}
	if err := a.store.DB.Create(&perm).Error; err != nil {
		jsonError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "权限已创建", "permission": makePermissionDTO(perm)})
}

func (a *App) apiPermissionsUpdate(c *gin.Context) {
	id, ok := parseAPIIDParam(c, "id")
	if !ok {
		return
	}
	var req permissionRequest
	if !bindJSON(c, &req) {
		return
	}
	var perm store.Permission
	if err := a.store.DB.First(&perm, id).Error; err != nil {
		jsonError(c, http.StatusNotFound, err)
		return
	}
	if req.ParentID == id {
		jsonMessage(c, http.StatusBadRequest, "上级权限不能选择自身")
		return
	}
	perm.ParentID = req.ParentID
	perm.MenuID = req.MenuID
	perm.Name = strings.TrimSpace(req.Name)
	perm.Code = strings.TrimSpace(req.Code)
	perm.Method = strings.ToUpper(strings.TrimSpace(req.Method))
	perm.Path = strings.TrimSpace(req.Path)
	perm.Sort = req.Sort
	perm.Status = req.Status
	if perm.Method == "" {
		perm.Method = http.MethodGet
	}
	if perm.Name == "" || perm.Code == "" {
		jsonMessage(c, http.StatusBadRequest, "权限名称和权限标识不能为空")
		return
	}
	if err := a.store.DB.Save(&perm).Error; err != nil {
		jsonError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "权限已更新", "permission": makePermissionDTO(perm)})
}

func (a *App) apiPermissionsDelete(c *gin.Context) {
	id, ok := parseAPIIDParam(c, "id")
	if !ok {
		return
	}
	err := a.store.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&store.Permission{}).Where("parent_id = ?", id).Update("parent_id", 0).Error; err != nil {
			return err
		}
		if err := tx.Where("permission_id = ?", id).Delete(&store.RolePermission{}).Error; err != nil {
			return err
		}
		return tx.Delete(&store.Permission{}, id).Error
	})
	if err != nil {
		jsonError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "权限已删除"})
}

func makeUserDTO(user store.User, roleIDs []uint) userDTO {
	roles := make([]roleDTO, 0, len(user.Roles))
	for _, role := range user.Roles {
		roles = append(roles, makeRoleDTO(role, nil, nil))
	}
	return userDTO{
		ID:       user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Email:    user.Email,
		Status:   user.Status,
		Roles:    roles,
		RoleIDs:  roleIDs,
	}
}

func makeRoleDTO(role store.Role, menuIDs, permIDs []uint) roleDTO {
	return roleDTO{
		ID:      role.ID,
		Name:    role.Name,
		Code:    role.Code,
		Remark:  role.Remark,
		Status:  role.Status,
		MenuIDs: menuIDs,
		PermIDs: permIDs,
	}
}

func makeMenuDTO(menu store.Menu) menuDTO {
	return menuDTO{
		ID:       menu.ID,
		ParentID: menu.ParentID,
		Title:    menu.Title,
		Icon:     menu.Icon,
		Path:     menu.Path,
		Sort:     menu.Sort,
		Visible:  menu.Visible,
		Status:   menu.Status,
	}
}

func makeMenuDTOs(menus []store.Menu) []menuDTO {
	out := make([]menuDTO, 0, len(menus))
	for _, menu := range menus {
		out = append(out, makeMenuDTO(menu))
	}
	return out
}

func makePermissionDTO(perm store.Permission) permissionDTO {
	return permissionDTO{
		ID:       perm.ID,
		ParentID: perm.ParentID,
		MenuID:   perm.MenuID,
		Name:     perm.Name,
		Code:     perm.Code,
		Method:   perm.Method,
		Path:     perm.Path,
		Sort:     perm.Sort,
		Status:   perm.Status,
		Menu:     makeMenuDTO(perm.Menu),
	}
}

func makePermissionDTOs(perms []store.Permission) []permissionDTO {
	out := make([]permissionDTO, 0, len(perms))
	for _, perm := range perms {
		out = append(out, makePermissionDTO(perm))
	}
	return out
}

func parseAPIIDParam(c *gin.Context, name string) (uint, bool) {
	value, err := strconv.ParseUint(c.Param(name), 10, 64)
	if err != nil || value == 0 {
		jsonMessage(c, http.StatusBadRequest, "无效 ID")
		return 0, false
	}
	return uint(value), true
}

func bindJSON(c *gin.Context, target any) bool {
	if err := c.ShouldBindJSON(target); err != nil {
		jsonError(c, http.StatusBadRequest, err)
		return false
	}
	return true
}

func jsonError(c *gin.Context, status int, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		jsonMessage(c, http.StatusNotFound, "记录不存在")
		return
	}
	jsonMessage(c, status, err.Error())
}

func jsonMessage(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}
