USE `shop`;

CREATE TABLE IF NOT EXISTS `site_settings` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `group_name` varchar(64) NOT NULL DEFAULT 'general',
  `label` varchar(128) NOT NULL,
  `key` varchar(96) NOT NULL,
  `value` text NOT NULL,
  `value_type` varchar(24) NOT NULL DEFAULT 'text',
  `sort` int NOT NULL DEFAULT 10,
  `status` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_site_settings_key` (`key`), KEY `idx_site_settings_group_sort` (`group_name`,`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `site_nav_items` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `parent_id` bigint unsigned NOT NULL DEFAULT 0, `label` varchar(128) NOT NULL,
  `path` varchar(500) NOT NULL DEFAULT '', `location` varchar(32) NOT NULL DEFAULT 'header', `open_new` tinyint(1) NOT NULL DEFAULT 0,
  `sort` int NOT NULL DEFAULT 10, `status` tinyint(1) NOT NULL DEFAULT 1, `created_at` datetime(3) DEFAULT NULL, `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`), KEY `idx_site_nav_items_parent_sort` (`parent_id`,`sort`), KEY `idx_site_nav_items_location` (`location`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `site_banners` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `title` varchar(255) NOT NULL DEFAULT '', `subtitle` varchar(500) NOT NULL DEFAULT '',
  `image_url` varchar(1000) NOT NULL DEFAULT '', `link_url` varchar(500) NOT NULL DEFAULT '', `position` varchar(32) NOT NULL DEFAULT 'home',
  `sort` int NOT NULL DEFAULT 10, `status` tinyint(1) NOT NULL DEFAULT 1, `created_at` datetime(3) DEFAULT NULL, `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`), KEY `idx_site_banners_position_sort` (`position`,`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `site_benefits` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `title` varchar(500) NOT NULL, `icon` varchar(64) NOT NULL DEFAULT 'check', `sort` int NOT NULL DEFAULT 10,
  `status` tinyint(1) NOT NULL DEFAULT 1, `created_at` datetime(3) DEFAULT NULL, `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`), KEY `idx_site_benefits_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `product_categories` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `parent_id` bigint unsigned NOT NULL DEFAULT 0, `name` varchar(160) NOT NULL, `slug` varchar(190) NOT NULL,
  `image_url` varchar(1000) NOT NULL DEFAULT '', `summary` text NOT NULL, `home_featured` tinyint(1) NOT NULL DEFAULT 0, `sort` int NOT NULL DEFAULT 10,
  `status` tinyint(1) NOT NULL DEFAULT 1, `created_at` datetime(3) DEFAULT NULL, `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_product_categories_slug` (`slug`), KEY `idx_product_categories_parent_sort` (`parent_id`,`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `products` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `category_id` bigint unsigned NOT NULL DEFAULT 0, `name` varchar(255) NOT NULL, `slug` varchar(255) NOT NULL,
  `sku` varchar(96) NOT NULL DEFAULT '', `summary` text NOT NULL, `description` longtext NOT NULL, `main_image` varchar(1000) NOT NULL DEFAULT '', `video_url` varchar(1000) NOT NULL DEFAULT '',
  `material` varchar(255) NOT NULL DEFAULT '', `size` varchar(255) NOT NULL DEFAULT '', `thread_standard` varchar(255) NOT NULL DEFAULT '',
  `pressure_rating` varchar(255) NOT NULL DEFAULT '', `temperature_range` varchar(255) NOT NULL DEFAULT '', `moq` varchar(255) NOT NULL DEFAULT '',
  `standard` varchar(255) NOT NULL DEFAULT '', `application` varchar(500) NOT NULL DEFAULT '', `hot_tags` varchar(1000) NOT NULL DEFAULT '',
  `featured` tinyint(1) NOT NULL DEFAULT 0, `latest` tinyint(1) NOT NULL DEFAULT 0, `sort` int NOT NULL DEFAULT 10, `status` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` datetime(3) DEFAULT NULL, `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_products_slug` (`slug`), KEY `idx_products_category_sort` (`category_id`,`sort`), KEY `idx_products_featured` (`featured`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `product_images` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `product_id` bigint unsigned NOT NULL, `image_url` varchar(1000) NOT NULL, `alt_text` varchar(255) NOT NULL DEFAULT '',
  `sort` int NOT NULL DEFAULT 10, `status` tinyint(1) NOT NULL DEFAULT 1, `created_at` datetime(3) DEFAULT NULL, `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`), KEY `idx_product_images_product_sort` (`product_id`,`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `content_pages` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `title` varchar(255) NOT NULL, `slug` varchar(190) NOT NULL, `subtitle` varchar(500) NOT NULL DEFAULT '',
  `content` longtext NOT NULL, `cover_image` varchar(1000) NOT NULL DEFAULT '', `template` varchar(32) NOT NULL DEFAULT 'standard', `sort` int NOT NULL DEFAULT 10,
  `status` tinyint(1) NOT NULL DEFAULT 1, `created_at` datetime(3) DEFAULT NULL, `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_content_pages_slug` (`slug`), KEY `idx_content_pages_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `content_articles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `content_type` varchar(32) NOT NULL, `category` varchar(128) NOT NULL DEFAULT '', `title` varchar(500) NOT NULL,
  `slug` varchar(255) NOT NULL, `summary` text NOT NULL, `content` longtext NOT NULL, `cover_image` varchar(1000) NOT NULL DEFAULT '',
  `published_at` datetime NOT NULL, `featured` tinyint(1) NOT NULL DEFAULT 0, `sort` int NOT NULL DEFAULT 10, `status` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` datetime(3) DEFAULT NULL, `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_content_articles_slug` (`slug`), KEY `idx_content_articles_type_date` (`content_type`,`published_at`), KEY `idx_content_articles_category` (`category`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `team_members` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `name` varchar(128) NOT NULL, `role` varchar(128) NOT NULL DEFAULT '', `email` varchar(190) NOT NULL DEFAULT '',
  `image_url` varchar(1000) NOT NULL DEFAULT '', `bio` text NOT NULL, `sort` int NOT NULL DEFAULT 10, `status` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` datetime(3) DEFAULT NULL, `updated_at` datetime(3) DEFAULT NULL, PRIMARY KEY (`id`), KEY `idx_team_members_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `partners` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `name` varchar(128) NOT NULL, `logo_url` varchar(1000) NOT NULL DEFAULT '', `website_url` varchar(1000) NOT NULL DEFAULT '',
  `sort` int NOT NULL DEFAULT 10, `status` tinyint(1) NOT NULL DEFAULT 1, `created_at` datetime(3) DEFAULT NULL, `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`), KEY `idx_partners_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `industries` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `title` varchar(190) NOT NULL, `subtitle` varchar(500) NOT NULL DEFAULT '', `description` text NOT NULL,
  `image_url` varchar(1000) NOT NULL DEFAULT '', `link_url` varchar(1000) NOT NULL DEFAULT '', `sort` int NOT NULL DEFAULT 10, `status` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` datetime(3) DEFAULT NULL, `updated_at` datetime(3) DEFAULT NULL, PRIMARY KEY (`id`), KEY `idx_industries_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `videos` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `title` varchar(255) NOT NULL, `description` text NOT NULL, `cover_url` varchar(1000) NOT NULL DEFAULT '',
  `video_url` varchar(1000) NOT NULL DEFAULT '', `sort` int NOT NULL DEFAULT 10, `status` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` datetime(3) DEFAULT NULL, `updated_at` datetime(3) DEFAULT NULL, PRIMARY KEY (`id`), KEY `idx_videos_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `inquiries` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `name` varchar(128) NOT NULL, `email` varchar(190) NOT NULL, `phone` varchar(128) NOT NULL DEFAULT '',
  `message` text NOT NULL, `attachment_url` varchar(1000) NOT NULL DEFAULT '', `product_ids` varchar(500) NOT NULL DEFAULT '', `source` varchar(96) NOT NULL DEFAULT '',
  `lead_status` varchar(32) NOT NULL DEFAULT 'new', `admin_note` text NOT NULL, `created_at` datetime(3) DEFAULT NULL, `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`), KEY `idx_inquiries_status_created` (`lead_status`,`created_at`), KEY `idx_inquiries_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `email_subscribers` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `email` varchar(190) NOT NULL, `status` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` datetime(3) DEFAULT NULL, `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_email_subscribers_email` (`email`), KEY `idx_email_subscribers_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT IGNORE INTO `admin_menus` (`id`,`parent_id`,`title`,`icon`,`path`,`sort`,`visible`,`status`,`created_at`,`updated_at`) VALUES
(10,0,'站点运营','W','',20,1,1,NOW(3),NOW(3)),
(11,10,'站点设置','S','/cms/settings',10,1,1,NOW(3),NOW(3)),
(12,10,'Banner 管理','B','/cms/banners',20,1,1,NOW(3),NOW(3)),
(13,10,'Why Depo','✓','/cms/benefits',30,1,1,NOW(3),NOW(3)),
(14,10,'合作伙伴','P','/cms/partners',40,1,1,NOW(3),NOW(3)),
(15,10,'团队成员','T','/cms/team',50,1,1,NOW(3),NOW(3)),
(16,10,'行业应用','I','/cms/industries',60,1,1,NOW(3),NOW(3)),
(18,10,'前台导航','☷','/cms/navigation',80,1,1,NOW(3),NOW(3)),
(20,0,'产品中心','◆','',30,1,1,NOW(3),NOW(3)),
(21,20,'产品分类','C','/cms/categories',10,1,1,NOW(3),NOW(3)),
(22,20,'产品管理','P','/cms/products',20,1,1,NOW(3),NOW(3)),
(30,0,'内容中心','N','',40,1,1,NOW(3),NOW(3)),
(31,30,'单页管理','D','/cms/pages',10,1,1,NOW(3),NOW(3)),
(32,30,'文章管理','A','/cms/articles',20,1,1,NOW(3),NOW(3)),
(33,30,'视频管理','V','/cms/videos',30,1,1,NOW(3),NOW(3)),
(40,0,'客户线索','L','',50,1,1,NOW(3),NOW(3)),
(41,40,'询盘管理','Q','/cms/inquiries',10,1,1,NOW(3),NOW(3)),
(42,40,'邮件订阅','E','/cms/subscribers',20,1,1,NOW(3),NOW(3));

INSERT IGNORE INTO `admin_permissions` (`id`,`parent_id`,`menu_id`,`name`,`code`,`method`,`path`,`sort`,`status`,`created_at`,`updated_at`) VALUES
(100,0,11,'站点设置查看','setting:view','GET','/api/cms/settings',10,1,NOW(3),NOW(3)),(101,100,11,'站点设置新增','setting:create','POST','/api/cms/settings',20,1,NOW(3),NOW(3)),(102,100,11,'站点设置编辑','setting:update','PUT','/api/cms/settings/:id',30,1,NOW(3),NOW(3)),(103,100,11,'站点设置删除','setting:delete','DELETE','/api/cms/settings/:id',40,1,NOW(3),NOW(3)),
(104,0,12,'Banner 查看','banner:view','GET','/api/cms/banners',10,1,NOW(3),NOW(3)),(105,104,12,'Banner 新增','banner:create','POST','/api/cms/banners',20,1,NOW(3),NOW(3)),(106,104,12,'Banner 编辑','banner:update','PUT','/api/cms/banners/:id',30,1,NOW(3),NOW(3)),(107,104,12,'Banner 删除','banner:delete','DELETE','/api/cms/banners/:id',40,1,NOW(3),NOW(3)),
(108,0,13,'Why Depo 查看','benefit:view','GET','/api/cms/benefits',10,1,NOW(3),NOW(3)),(109,108,13,'Why Depo 新增','benefit:create','POST','/api/cms/benefits',20,1,NOW(3),NOW(3)),(110,108,13,'Why Depo 编辑','benefit:update','PUT','/api/cms/benefits/:id',30,1,NOW(3),NOW(3)),(111,108,13,'Why Depo 删除','benefit:delete','DELETE','/api/cms/benefits/:id',40,1,NOW(3),NOW(3)),
(112,0,21,'产品分类查看','category:view','GET','/api/cms/categories',10,1,NOW(3),NOW(3)),(113,112,21,'产品分类新增','category:create','POST','/api/cms/categories',20,1,NOW(3),NOW(3)),(114,112,21,'产品分类编辑','category:update','PUT','/api/cms/categories/:id',30,1,NOW(3),NOW(3)),(115,112,21,'产品分类删除','category:delete','DELETE','/api/cms/categories/:id',40,1,NOW(3),NOW(3)),
(116,0,22,'产品查看','product:view','GET','/api/cms/products',10,1,NOW(3),NOW(3)),(117,116,22,'产品新增','product:create','POST','/api/cms/products',20,1,NOW(3),NOW(3)),(118,116,22,'产品编辑','product:update','PUT','/api/cms/products/:id',30,1,NOW(3),NOW(3)),(119,116,22,'产品删除','product:delete','DELETE','/api/cms/products/:id',40,1,NOW(3),NOW(3)),
(124,0,31,'单页查看','page:view','GET','/api/cms/pages',10,1,NOW(3),NOW(3)),(125,124,31,'单页新增','page:create','POST','/api/cms/pages',20,1,NOW(3),NOW(3)),(126,124,31,'单页编辑','page:update','PUT','/api/cms/pages/:id',30,1,NOW(3),NOW(3)),(127,124,31,'单页删除','page:delete','DELETE','/api/cms/pages/:id',40,1,NOW(3),NOW(3)),
(128,0,32,'文章查看','article:view','GET','/api/cms/articles',10,1,NOW(3),NOW(3)),(129,128,32,'文章新增','article:create','POST','/api/cms/articles',20,1,NOW(3),NOW(3)),(130,128,32,'文章编辑','article:update','PUT','/api/cms/articles/:id',30,1,NOW(3),NOW(3)),(131,128,32,'文章删除','article:delete','DELETE','/api/cms/articles/:id',40,1,NOW(3),NOW(3)),
(132,0,15,'团队查看','team:view','GET','/api/cms/team',10,1,NOW(3),NOW(3)),(133,132,15,'团队新增','team:create','POST','/api/cms/team',20,1,NOW(3),NOW(3)),(134,132,15,'团队编辑','team:update','PUT','/api/cms/team/:id',30,1,NOW(3),NOW(3)),(135,132,15,'团队删除','team:delete','DELETE','/api/cms/team/:id',40,1,NOW(3),NOW(3)),
(136,0,14,'合作伙伴查看','partner:view','GET','/api/cms/partners',10,1,NOW(3),NOW(3)),(137,136,14,'合作伙伴新增','partner:create','POST','/api/cms/partners',20,1,NOW(3),NOW(3)),(138,136,14,'合作伙伴编辑','partner:update','PUT','/api/cms/partners/:id',30,1,NOW(3),NOW(3)),(139,136,14,'合作伙伴删除','partner:delete','DELETE','/api/cms/partners/:id',40,1,NOW(3),NOW(3)),
(140,0,16,'行业应用查看','industry:view','GET','/api/cms/industries',10,1,NOW(3),NOW(3)),(141,140,16,'行业应用新增','industry:create','POST','/api/cms/industries',20,1,NOW(3),NOW(3)),(142,140,16,'行业应用编辑','industry:update','PUT','/api/cms/industries/:id',30,1,NOW(3),NOW(3)),(143,140,16,'行业应用删除','industry:delete','DELETE','/api/cms/industries/:id',40,1,NOW(3),NOW(3)),
(144,0,33,'视频查看','video:view','GET','/api/cms/videos',10,1,NOW(3),NOW(3)),(145,144,33,'视频新增','video:create','POST','/api/cms/videos',20,1,NOW(3),NOW(3)),(146,144,33,'视频编辑','video:update','PUT','/api/cms/videos/:id',30,1,NOW(3),NOW(3)),(147,144,33,'视频删除','video:delete','DELETE','/api/cms/videos/:id',40,1,NOW(3),NOW(3)),
(152,0,41,'询盘查看','inquiry:view','GET','/api/cms/inquiries',10,1,NOW(3),NOW(3)),(153,152,41,'询盘跟进','inquiry:update','PUT','/api/cms/inquiries/:id',20,1,NOW(3),NOW(3)),(154,152,41,'询盘删除','inquiry:delete','DELETE','/api/cms/inquiries/:id',30,1,NOW(3),NOW(3)),
(155,0,42,'订阅查看','subscriber:view','GET','/api/cms/subscribers',10,1,NOW(3),NOW(3)),(156,155,42,'订阅新增','subscriber:create','POST','/api/cms/subscribers',20,1,NOW(3),NOW(3)),(157,155,42,'订阅编辑','subscriber:update','PUT','/api/cms/subscribers/:id',30,1,NOW(3),NOW(3)),(158,155,42,'订阅删除','subscriber:delete','DELETE','/api/cms/subscribers/:id',40,1,NOW(3),NOW(3));

INSERT IGNORE INTO `admin_permissions` (`id`,`parent_id`,`menu_id`,`name`,`code`,`method`,`path`,`sort`,`status`,`created_at`,`updated_at`) VALUES
(159,0,18,'前台导航查看','navigation:view','GET','/api/cms/navigation',10,1,NOW(3),NOW(3)),
(160,159,18,'前台导航新增','navigation:create','POST','/api/cms/navigation',20,1,NOW(3),NOW(3)),
(161,159,18,'前台导航编辑','navigation:update','PUT','/api/cms/navigation/:id',30,1,NOW(3),NOW(3)),
(162,159,18,'前台导航删除','navigation:delete','DELETE','/api/cms/navigation/:id',40,1,NOW(3),NOW(3));

INSERT IGNORE INTO `admin_role_menus` (`role_id`,`menu_id`) SELECT 1, id FROM `admin_menus` WHERE id >= 10;
INSERT IGNORE INTO `admin_role_permissions` (`role_id`,`permission_id`) SELECT 1, id FROM `admin_permissions` WHERE id >= 100;
INSERT IGNORE INTO `admin_role_menus` (`role_id`,`menu_id`) SELECT 2, id FROM `admin_menus` WHERE id >= 10;
INSERT IGNORE INTO `admin_role_permissions` (`role_id`,`permission_id`) SELECT 2, id FROM `admin_permissions` WHERE id IN (100,104,108,112,116,124,128,132,136,140,144,152,155,159);

-- 清理已初始化数据库中的旧图片目录和独立产品图集入口；业务数据表及历史图片数据保留。
DELETE `arp` FROM `admin_role_permissions` AS `arp`
JOIN `admin_permissions` AS `ap` ON `ap`.`id` = `arp`.`permission_id`
WHERE `ap`.`menu_id` IN (17,23);
DELETE FROM `admin_role_menus` WHERE `menu_id` IN (17,23);
DELETE FROM `admin_permissions` WHERE `menu_id` IN (17,23);
DELETE FROM `admin_menus` WHERE `id` IN (17,23);
