# 项目部署文档

本文档记录当前服务器部署结构、服务配置位置、常用启动关闭命令，以及后续更新发布步骤。

## 当前部署信息

- 服务器：`8.166.143.93`
- 前台地址：`http://8.166.143.93/`
- 后台地址：`http://8.166.143.93:8080/login`
- 前台 API：`http://8.166.143.93/api/v1/`
- 发布目录：`/opt/shop`
- MySQL 数据库：`shop`
- MySQL 用户：`shop_admin`

当前未使用域名，Nginx 直接按服务器 IP 提供访问。

## 远程目录结构

```text
/opt/shop
├── admin
│   ├── rbac-admin
│   ├── .env
│   └── web
│       ├── static
│       └── templates
└── front
    ├── hbfittings-front-api
    ├── .env
    └── web
        └── dist
```

上传目录：

```text
/opt/shop/admin/web/static/uploads
```

前台和后台共用这个上传目录。Nginx 通过 `/uploads/` 代理到前台 API。

## 服务说明

当前有两个 systemd 服务：

```text
shop-admin.service
shop-front-api.service
```

服务监听方式：

```text
后台 Go 服务      127.0.0.1:18080
前台 API 服务     127.0.0.1:8090
Nginx 前台入口    0.0.0.0:80
Nginx 后台入口    0.0.0.0:8080
```

外部访问不直接访问 Go 服务，由 Nginx 反向代理。

## 环境变量配置

后台配置文件：

```bash
/opt/shop/admin/.env
```

示例：

```bash
RBAC_ADDR=127.0.0.1:18080
RBAC_DSN="shop_admin:<mysql_password>@tcp(127.0.0.1:3306)/shop?charset=utf8mb4&parseTime=True&loc=Local"
RBAC_SESSION_SECRET="<random_secret>"
RBAC_UPLOAD_DIR=/opt/shop/admin/web/static/uploads
```

前台 API 配置文件：

```bash
/opt/shop/front/.env
```

示例：

```bash
FRONT_API_ADDR=127.0.0.1:8090
FRONT_DSN="shop_admin:<mysql_password>@tcp(127.0.0.1:3306)/shop?charset=utf8mb4&parseTime=True&loc=Local"
FRONT_UPLOAD_DIR=/opt/shop/admin/web/static/uploads
```

不要把真实数据库密码和 session secret 提交到代码仓库。

## Nginx 配置

当前 Nginx 站点配置文件：

```bash
/www/server/panel/vhost/nginx/8.166.143.93.conf
```

核心代理关系：

```text
http://8.166.143.93/              -> /opt/shop/front/web/dist
http://8.166.143.93/api/...       -> 127.0.0.1:8090
http://8.166.143.93/uploads/...   -> 127.0.0.1:8090
http://8.166.143.93:8080/...      -> 127.0.0.1:18080
```

检查并重载 Nginx：

```bash
nginx -t
systemctl reload nginx
```

重启 Nginx：

```bash
systemctl restart nginx
```

## 启动、停止、重启

查看服务状态：

```bash
systemctl status shop-admin.service --no-pager -l
systemctl status shop-front-api.service --no-pager -l
systemctl status nginx --no-pager -l
```

启动服务：

```bash
systemctl start shop-admin.service
systemctl start shop-front-api.service
systemctl start nginx
```

停止服务：

```bash
systemctl stop shop-admin.service
systemctl stop shop-front-api.service
systemctl stop nginx
```

重启服务：

```bash
systemctl restart shop-admin.service
systemctl restart shop-front-api.service
systemctl restart nginx
```

设置开机自启：

```bash
systemctl enable shop-admin.service
systemctl enable shop-front-api.service
systemctl enable nginx
```

取消开机自启：

```bash
systemctl disable shop-admin.service
systemctl disable shop-front-api.service
```

## 日志查看

后台服务日志：

```bash
journalctl -u shop-admin.service -n 100 --no-pager
journalctl -u shop-admin.service -f
```

前台 API 日志：

```bash
journalctl -u shop-front-api.service -n 100 --no-pager
journalctl -u shop-front-api.service -f
```

Nginx 日志：

```bash
tail -100 /www/wwwlogs/shop-front.log
tail -100 /www/wwwlogs/shop-front.error.log
tail -100 /www/wwwlogs/shop-admin.log
tail -100 /www/wwwlogs/shop-admin.error.log
```

## 本地构建

在本地项目根目录执行。

构建前台静态资源：

```bash
cd front/web
npm run build
```

构建 Linux amd64 后台二进制：

```bash
cd admin
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /tmp/rbac-admin ./cmd/server
```

构建 Linux amd64 前台 API 二进制：

```bash
cd front
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /tmp/hbfittings-front-api ./cmd/server
```

## 发布更新

只更新前台页面：

```bash
cd front/web
npm run build
rsync -az dist/ root@8.166.143.93:/opt/shop/front/web/dist/
```

只更新前台 API：

```bash
cd front
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /tmp/hbfittings-front-api ./cmd/server
scp /tmp/hbfittings-front-api root@8.166.143.93:/opt/shop/front/hbfittings-front-api
ssh root@8.166.143.93 'chmod +x /opt/shop/front/hbfittings-front-api && systemctl restart shop-front-api.service'
```

只更新后台：

```bash
cd admin
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /tmp/rbac-admin ./cmd/server
scp /tmp/rbac-admin root@8.166.143.93:/opt/shop/admin/rbac-admin
rsync -az web/ root@8.166.143.93:/opt/shop/admin/web/
ssh root@8.166.143.93 'chmod +x /opt/shop/admin/rbac-admin && systemctl restart shop-admin.service'
```

完整更新：

```bash
# 1. 构建前台页面
cd front/web
npm run build

# 2. 构建前台 API
cd ../
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /tmp/hbfittings-front-api ./cmd/server

# 3. 构建后台
cd ../admin
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /tmp/rbac-admin ./cmd/server

# 4. 上传
scp /tmp/hbfittings-front-api root@8.166.143.93:/opt/shop/front/hbfittings-front-api
scp /tmp/rbac-admin root@8.166.143.93:/opt/shop/admin/rbac-admin
rsync -az ../front/web/dist/ root@8.166.143.93:/opt/shop/front/web/dist/
rsync -az web/ root@8.166.143.93:/opt/shop/admin/web/

# 5. 重启服务
ssh root@8.166.143.93 '
chmod +x /opt/shop/front/hbfittings-front-api /opt/shop/admin/rbac-admin
systemctl restart shop-front-api.service shop-admin.service
nginx -t && systemctl reload nginx
'
```

## 数据库注意事项

远程 `shop` 数据库已经初始化并有数据。不要随意执行：

```bash
admin/front_seed.sql
```

该 seed 文件包含 `DELETE FROM ...`，会清空演示内容相关表后重新导入数据。

如果只是补表结构，可先检查 SQL 内容，确认不会删除线上数据后再执行。

## 访问验证

前台：

```bash
curl -I http://8.166.143.93/
```

前台 API：

```bash
curl -I http://8.166.143.93/api/v1/bootstrap
curl http://8.166.143.93/api/v1/bootstrap | head -c 300
```

后台：

```bash
curl -I http://8.166.143.93:8080/login
```

远程端口监听：

```bash
ss -lntp | egrep ':80|:8080|:18080|:8090'
```

正常状态应包含：

```text
0.0.0.0:80       nginx
0.0.0.0:8080     nginx
127.0.0.1:18080  rbac-admin
127.0.0.1:8090   hbfittings-front-api
```

## 防火墙端口

服务器当前需要开放：

```text
80/tcp
8080/tcp
22/tcp
```

查看防火墙端口：

```bash
firewall-cmd --list-ports
```

放行后台端口：

```bash
firewall-cmd --add-port=8080/tcp --permanent
firewall-cmd --reload
```

如果公网仍无法访问 `8080`，还需要在云服务器安全组中放行 `8080/tcp`。
