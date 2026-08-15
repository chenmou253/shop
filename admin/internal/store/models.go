package store

import "time"

type User struct {
	ID           uint `gorm:"primaryKey"`
	Username     string
	PasswordHash string
	PasswordSalt string
	Nickname     string
	Email        string
	Status       bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Roles        []Role `gorm:"many2many:admin_user_roles;"`
}

func (User) TableName() string { return "admin_users" }

type Role struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	Code      string
	Remark    string
	Status    bool
	CreatedAt time.Time
	UpdatedAt time.Time
	Users     []User       `gorm:"many2many:admin_user_roles;"`
	Menus     []Menu       `gorm:"many2many:admin_role_menus;"`
	Perms     []Permission `gorm:"many2many:admin_role_permissions;"`
}

func (Role) TableName() string { return "admin_roles" }

type Menu struct {
	ID        uint `gorm:"primaryKey"`
	ParentID  uint
	Title     string
	Icon      string
	Path      string
	Sort      int
	Visible   bool
	Status    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Menu) TableName() string { return "admin_menus" }

type Permission struct {
	ID        uint `gorm:"primaryKey"`
	ParentID  uint
	MenuID    uint
	Name      string
	Code      string
	Method    string
	Path      string
	Sort      int
	Status    bool
	CreatedAt time.Time
	UpdatedAt time.Time
	Menu      Menu `gorm:"foreignKey:MenuID"`
}

func (Permission) TableName() string { return "admin_permissions" }

type UserRole struct {
	UserID uint `gorm:"primaryKey"`
	RoleID uint `gorm:"primaryKey"`
}

func (UserRole) TableName() string { return "admin_user_roles" }

type RoleMenu struct {
	RoleID uint `gorm:"primaryKey"`
	MenuID uint `gorm:"primaryKey"`
}

func (RoleMenu) TableName() string { return "admin_role_menus" }

type RolePermission struct {
	RoleID       uint `gorm:"primaryKey"`
	PermissionID uint `gorm:"primaryKey"`
}

func (RolePermission) TableName() string { return "admin_role_permissions" }
