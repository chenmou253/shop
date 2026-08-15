package store

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"sort"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Store struct {
	DB *gorm.DB
}

func Open(dsn string) (*Store, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &Store{DB: db}, nil
}

func HashPassword(salt, password string) string {
	sum := sha256.Sum256([]byte(salt + ":" + password))
	return hex.EncodeToString(sum[:])
}

func RandomSalt() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func (s *Store) Authenticate(username, password string) (*User, error) {
	var user User
	if err := s.DB.Preload("Roles").Where("username = ? AND status = ?", username, true).First(&user).Error; err != nil {
		return nil, err
	}
	if HashPassword(user.PasswordSalt, password) != user.PasswordHash {
		return nil, errors.New("invalid password")
	}
	return &user, nil
}

func (s *Store) UserHasPermission(userID uint, code string) bool {
	if code == "" {
		return true
	}
	var count int64
	s.DB.Table("admin_permissions p").
		Joins("JOIN admin_role_permissions rp ON rp.permission_id = p.id").
		Joins("JOIN admin_user_roles ur ON ur.role_id = rp.role_id").
		Joins("JOIN admin_roles r ON r.id = ur.role_id").
		Where("ur.user_id = ? AND p.code = ? AND p.status = ? AND r.status = ?", userID, code, true, true).
		Count(&count)
	return count > 0
}

func (s *Store) RoleIDsForUser(userID uint) ([]uint, error) {
	var rows []UserRole
	if err := s.DB.Where("user_id = ?", userID).Find(&rows).Error; err != nil {
		return nil, err
	}
	ids := make([]uint, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.RoleID)
	}
	return ids, nil
}

func (s *Store) MenuIDsForRole(roleID uint) ([]uint, error) {
	var rows []RoleMenu
	if err := s.DB.Where("role_id = ?", roleID).Find(&rows).Error; err != nil {
		return nil, err
	}
	ids := make([]uint, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.MenuID)
	}
	return ids, nil
}

func (s *Store) PermissionIDsForRole(roleID uint) ([]uint, error) {
	var rows []RolePermission
	if err := s.DB.Where("role_id = ?", roleID).Find(&rows).Error; err != nil {
		return nil, err
	}
	ids := make([]uint, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.PermissionID)
	}
	return ids, nil
}

func (s *Store) ReplaceUserRoles(tx *gorm.DB, userID uint, roleIDs []uint) error {
	if err := tx.Where("user_id = ?", userID).Delete(&UserRole{}).Error; err != nil {
		return err
	}
	rows := make([]UserRole, 0, len(roleIDs))
	for _, roleID := range uniqueUint(roleIDs) {
		rows = append(rows, UserRole{UserID: userID, RoleID: roleID})
	}
	if len(rows) == 0 {
		return nil
	}
	return tx.Create(&rows).Error
}

func (s *Store) ReplaceRoleAccess(tx *gorm.DB, roleID uint, menuIDs, permissionIDs []uint) error {
	if err := tx.Where("role_id = ?", roleID).Delete(&RoleMenu{}).Error; err != nil {
		return err
	}
	if err := tx.Where("role_id = ?", roleID).Delete(&RolePermission{}).Error; err != nil {
		return err
	}

	menuRows := make([]RoleMenu, 0, len(menuIDs))
	for _, menuID := range uniqueUint(menuIDs) {
		menuRows = append(menuRows, RoleMenu{RoleID: roleID, MenuID: menuID})
	}
	if len(menuRows) > 0 {
		if err := tx.Create(&menuRows).Error; err != nil {
			return err
		}
	}

	permRows := make([]RolePermission, 0, len(permissionIDs))
	for _, permissionID := range uniqueUint(permissionIDs) {
		permRows = append(permRows, RolePermission{RoleID: roleID, PermissionID: permissionID})
	}
	if len(permRows) > 0 {
		if err := tx.Create(&permRows).Error; err != nil {
			return err
		}
	}
	return nil
}

func ParseUintList(values []string) []uint {
	ids := make([]uint, 0, len(values))
	for _, value := range values {
		n, err := strconv.ParseUint(value, 10, 64)
		if err != nil || n == 0 {
			continue
		}
		ids = append(ids, uint(n))
	}
	return uniqueUint(ids)
}

func uniqueUint(ids []uint) []uint {
	seen := make(map[uint]struct{}, len(ids))
	out := make([]uint, 0, len(ids))
	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}
