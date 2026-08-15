# Depo Website Admin

一个 Gin + GORM 的前台内容管理与 RBAC 后台。除用户、角色、菜单和权限外，还完整管理 `front` 使用的站点设置、Banner、产品、内容和客户线索。各业务新增、编辑表单内可直接上传图片，产品图集也统一在产品表单中维护。

当前实现为页面 + API 分离：

- `/login`、`/dashboard`、`/users`、`/roles`、`/menus`、`/permissions` 只负责返回页面壳。
- 登录、列表、详情、新增、编辑、删除、树形联动全部通过 `/api/*` JSON 接口完成。
- 左侧菜单也通过 `/api/nav` 按当前登录用户权限加载。

## 初始化

1. 创建并初始化基础 RBAC 数据库：

```bash
mysql -uroot -p < init.sql
```

2. 创建前台共享内容表、后台菜单和权限：

```bash
mysql -uroot -p < front_schema.sql
```

3. 导入英文演示数据（所有默认图片均使用 `https://www.hbfittings.net/` 资源）：

```bash
mysql -uroot -p < front_seed.sql
```

4. 按需设置数据库连接和图片上传目录：

```bash
export RBAC_DSN='root:password@tcp(127.0.0.1:3306)/rbac_admin?charset=utf8mb4&parseTime=True&loc=Local'
export RBAC_SESSION_SECRET='replace-with-a-random-secret'
export RBAC_UPLOAD_DIR='web/static/uploads'
```

5. 启动：

```bash
go run ./cmd/server
```

默认地址：`http://127.0.0.1:8080`。上传图片通过 `/uploads/...` 访问，`front` API 默认读取同一目录。

默认账号：`admin`

默认密码：`admin123`

## 目录

```text
cmd/server              启动入口
internal/app            Gin 路由、处理器、页面数据
internal/session        Cookie 会话
internal/store          GORM 模型、密码和关联表操作
web/templates           页面模板
web/static              CSS 和 JS
init.sql                手动执行的数据库初始化脚本
front_schema.sql         front 共享表、后台菜单与权限
front_seed.sql           front 英文动态内容演示数据
```

## API 概览

```text
POST   /api/login
POST   /api/logout
GET    /api/me
GET    /api/nav
GET    /api/dashboard

GET    /api/users
POST   /api/users
GET    /api/users/:id
PUT    /api/users/:id
DELETE /api/users/:id

GET    /api/roles
POST   /api/roles
GET    /api/roles/:id
PUT    /api/roles/:id
DELETE /api/roles/:id

GET    /api/menus
GET    /api/menus/tree
POST   /api/menus
GET    /api/menus/:id
PUT    /api/menus/:id
DELETE /api/menus/:id

GET    /api/permissions
GET    /api/permissions/tree
POST   /api/permissions
GET    /api/permissions/:id
PUT    /api/permissions/:id
DELETE /api/permissions/:id

GET    /api/role-access
GET    /api/permission-parents

GET    /api/cms/:resource/meta
GET    /api/cms/:resource
POST   /api/cms/:resource
GET    /api/cms/:resource/:id
PUT    /api/cms/:resource/:id
DELETE /api/cms/:resource/:id
POST   /api/cms/:resource/upload
```

## 前后台数据映射

`front` 与 `admin` 使用同一个 `shop` 数据库。前台导航、首页轮播、卖点、伙伴、团队、行业应用、产品分类/详情、News/Knowledge/Blog、单页、视频、联系方式、询盘和邮件订阅均有对应后台菜单，并按 `view/create/update/delete` 拆分 RBAC 权限。图片在相应业务的新增、编辑表单内上传；产品主图和产品图集统一在产品管理中维护。询盘由前台提交，后台负责查看、跟进和删除。
