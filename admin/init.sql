CREATE DATABASE IF NOT EXISTS `shop` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `shop`;

SET FOREIGN_KEY_CHECKS = 0;
DROP TABLE IF EXISTS `admin_role_permissions`;
DROP TABLE IF EXISTS `admin_role_menus`;
DROP TABLE IF EXISTS `admin_user_roles`;
DROP TABLE IF EXISTS `admin_permissions`;
DROP TABLE IF EXISTS `admin_menus`;
DROP TABLE IF EXISTS `admin_roles`;
DROP TABLE IF EXISTS `admin_users`;
SET FOREIGN_KEY_CHECKS = 1;

CREATE TABLE `admin_users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(64) NOT NULL,
  `password_hash` char(64) NOT NULL,
  `password_salt` varchar(64) NOT NULL,
  `nickname` varchar(64) NOT NULL DEFAULT '',
  `email` varchar(128) NOT NULL DEFAULT '',
  `status` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_admin_users_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `admin_roles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL,
  `code` varchar(64) NOT NULL,
  `remark` varchar(255) NOT NULL DEFAULT '',
  `status` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_admin_roles_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `admin_menus` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `parent_id` bigint unsigned NOT NULL DEFAULT 0,
  `title` varchar(64) NOT NULL,
  `icon` varchar(32) NOT NULL DEFAULT '',
  `path` varchar(128) NOT NULL DEFAULT '',
  `sort` int NOT NULL DEFAULT 10,
  `visible` tinyint(1) NOT NULL DEFAULT 1,
  `status` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_admin_menus_parent_id` (`parent_id`),
  KEY `idx_admin_menus_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `admin_permissions` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `parent_id` bigint unsigned NOT NULL DEFAULT 0,
  `menu_id` bigint unsigned NOT NULL DEFAULT 0,
  `name` varchar(64) NOT NULL,
  `code` varchar(96) NOT NULL,
  `method` varchar(12) NOT NULL DEFAULT 'GET',
  `path` varchar(160) NOT NULL DEFAULT '',
  `sort` int NOT NULL DEFAULT 10,
  `status` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_admin_permissions_code` (`code`),
  KEY `idx_admin_permissions_parent_id` (`parent_id`),
  KEY `idx_admin_permissions_menu_id` (`menu_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `admin_user_roles` (
  `user_id` bigint unsigned NOT NULL,
  `role_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`user_id`, `role_id`),
  KEY `idx_admin_user_roles_role_id` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `admin_role_menus` (
  `role_id` bigint unsigned NOT NULL,
  `menu_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`role_id`, `menu_id`),
  KEY `idx_admin_role_menus_menu_id` (`menu_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `admin_role_permissions` (
  `role_id` bigint unsigned NOT NULL,
  `permission_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`role_id`, `permission_id`),
  KEY `idx_admin_role_permissions_permission_id` (`permission_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO `admin_users` (`id`, `username`, `password_hash`, `password_salt`, `nickname`, `email`, `status`, `created_at`, `updated_at`) VALUES
(1, 'admin', '4db7889618b5c81e80e41839f01a8cdec1060a6293fcf20868c0574610e029ce', 'rbac-admin-salt', '超级管理员', 'admin@example.com', 1, NOW(3), NOW(3));

INSERT INTO `admin_roles` (`id`, `name`, `code`, `remark`, `status`, `created_at`, `updated_at`) VALUES
(1, '超级管理员', 'super_admin', '系统内置最高权限角色', 1, NOW(3), NOW(3)),
(2, '只读管理员', 'viewer', '仅查看基础后台数据', 1, NOW(3), NOW(3));

INSERT INTO `admin_menus` (`id`, `parent_id`, `title`, `icon`, `path`, `sort`, `visible`, `status`, `created_at`, `updated_at`) VALUES
(1, 0, '系统管理', '▦', '', 10, 1, 1, NOW(3), NOW(3)),
(2, 1, '用户管理', 'U', '/users', 20, 1, 1, NOW(3), NOW(3)),
(3, 1, '角色管理', 'R', '/roles', 30, 1, 1, NOW(3), NOW(3)),
(4, 1, '菜单管理', 'M', '/menus', 40, 1, 1, NOW(3), NOW(3)),
(5, 1, '权限管理', 'P', '/permissions', 50, 1, 1, NOW(3), NOW(3));

INSERT INTO `admin_permissions` (`id`, `parent_id`, `menu_id`, `name`, `code`, `method`, `path`, `sort`, `status`, `created_at`, `updated_at`) VALUES
(1, 0, 2, '用户查看', 'user:view', 'GET', '/users', 10, 1, NOW(3), NOW(3)),
(2, 1, 2, '用户新增', 'user:create', 'POST', '/users', 20, 1, NOW(3), NOW(3)),
(3, 1, 2, '用户编辑', 'user:update', 'POST', '/users/:id', 30, 1, NOW(3), NOW(3)),
(4, 1, 2, '用户删除', 'user:delete', 'POST', '/users/:id/delete', 40, 1, NOW(3), NOW(3)),
(5, 0, 3, '角色查看', 'role:view', 'GET', '/roles', 10, 1, NOW(3), NOW(3)),
(6, 5, 3, '角色新增', 'role:create', 'POST', '/roles', 20, 1, NOW(3), NOW(3)),
(7, 5, 3, '角色编辑', 'role:update', 'POST', '/roles/:id', 30, 1, NOW(3), NOW(3)),
(8, 5, 3, '角色删除', 'role:delete', 'POST', '/roles/:id/delete', 40, 1, NOW(3), NOW(3)),
(9, 0, 4, '菜单查看', 'menu:view', 'GET', '/menus', 10, 1, NOW(3), NOW(3)),
(10, 9, 4, '菜单新增', 'menu:create', 'POST', '/menus', 20, 1, NOW(3), NOW(3)),
(11, 9, 4, '菜单编辑', 'menu:update', 'POST', '/menus/:id', 30, 1, NOW(3), NOW(3)),
(12, 9, 4, '菜单删除', 'menu:delete', 'POST', '/menus/:id/delete', 40, 1, NOW(3), NOW(3)),
(13, 0, 5, '权限查看', 'permission:view', 'GET', '/permissions', 10, 1, NOW(3), NOW(3)),
(14, 13, 5, '权限新增', 'permission:create', 'POST', '/permissions', 20, 1, NOW(3), NOW(3)),
(15, 13, 5, '权限编辑', 'permission:update', 'POST', '/permissions/:id', 30, 1, NOW(3), NOW(3)),
(16, 13, 5, '权限删除', 'permission:delete', 'POST', '/permissions/:id/delete', 40, 1, NOW(3), NOW(3));

INSERT INTO `admin_user_roles` (`user_id`, `role_id`) VALUES
(1, 1);

INSERT INTO `admin_role_menus` (`role_id`, `menu_id`) VALUES
(1, 1), (1, 2), (1, 3), (1, 4), (1, 5),
(2, 1), (2, 2), (2, 3), (2, 4), (2, 5);

INSERT INTO `admin_role_permissions` (`role_id`, `permission_id`) VALUES
(1, 1), (1, 2), (1, 3), (1, 4), (1, 5), (1, 6), (1, 7), (1, 8), (1, 9), (1, 10), (1, 11), (1, 12), (1, 13), (1, 14), (1, 15), (1, 16),
(2, 1), (2, 5), (2, 9), (2, 13);
