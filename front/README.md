# Depo Front

`front` 是 hbfittings.net 英文站的前后端分离实现：

- `web/`：React + TypeScript + Vite。
- `cmd/server`、`internal/`：Gin + GORM 内容 API。
- 数据库：与 `../admin` 共用 MySQL `shop` 数据库。
- 图片：初始化数据使用原站 `https://www.hbfittings.net/` 地址；后台上传图片统一写入 `../admin/web/static/uploads`，两个服务都通过 `/uploads/...` 访问。

## 初始化数据库

```bash
mysql -uroot -p < ../admin/init.sql
mysql -uroot -p < ../admin/front_schema.sql
mysql -uroot -p < ../admin/front_seed.sql
```

## 启动 API

```bash
export FRONT_DSN='root:password@tcp(127.0.0.1:3306)/shop?charset=utf8mb4&parseTime=True&loc=Local'
export FRONT_API_ADDR=':8090'
export FRONT_UPLOAD_DIR='../admin/web/static/uploads'
go run ./cmd/server
```

## 启动 React

```bash
cd web
npm install
npm run dev
```

浏览器打开 `http://127.0.0.1:5173`。Vite 会把 `/api` 和 `/uploads` 代理到 `http://127.0.0.1:8090`。

生产构建：

```bash
cd web
npm run build
```

如 API 与网站不同域部署，可在构建时设置 `VITE_API_BASE=https://api.example.com/api/v1`。
生产环境还需要把网站的 `/uploads/` 反向代理到 Gin API 的 `/uploads/`，以便显示 admin 上传的图片。
