# 个人博客

一个面向个人使用的轻量博客系统。访客可以阅读 Markdown 文章和提交评论；管理员可以在独立后台发布文章、审核评论并查看内容统计。

项目使用前后端分离架构，但通过 Docker Compose 和 Nginx 对外提供单一入口：

- 博客首页：`http://localhost:8080`
- 管理后台：`http://localhost:8080/admin/login`
- 健康检查：`http://localhost:8080/api/v1/health`

## 技术栈

- 前端：Vue 3、TypeScript、Vite、Vue Router、Pinia、Axios、Element Plus
- 内容渲染：markdown-it、DOMPurify
- 后端：Go、Gin、GORM、JWT、bcrypt、`log/slog`
- 数据库：PostgreSQL 17
- 部署：Docker、Docker Compose、Nginx

## 目录结构

```text
.
├── backend/                Go API 服务
│   ├── cmd/server/         服务入口
│   └── internal/           配置、模型、中间件、Handler 与 Service
├── frontend/               Vue 前端
│   └── src/                API、组件、布局、路由、状态与页面
├── deploy/nginx.conf       API 反向代理与 Vue History 回退
├── docs/REQUIREMENTS.md    项目需求与验收标准
├── .env.example            环境变量模板
├── docker-compose.yml      三服务编排
└── README.md
```

## 环境要求

推荐直接使用 Docker：

- Docker Engine 或 Docker Desktop
- Docker Compose v2（使用 `docker compose` 命令）

如果需要脱离 Docker 开发，还需要：

- Go 1.24+
- Node.js 22+
- npm
- PostgreSQL 17 或兼容版本

## Docker 快速启动

1. 创建本地环境变量文件：

   ```bash
   cp .env.example .env
   ```

2. 编辑 `.env`，至少将以下值替换为自己的安全值，并确保数据库密码在 `POSTGRES_PASSWORD` 和 `DATABASE_URL` 中保持一致：

   ```env
   POSTGRES_PASSWORD=请设置数据库密码
   DATABASE_URL=host=postgres user=blog_user password=请设置数据库密码 dbname=personal_blog port=5432 sslmode=disable
   ADMIN_PASSWORD=请设置管理员密码
   JWT_SECRET=请设置至少32个字符的随机字符串
   ```

3. 构建并启动全部服务：

   ```bash
   docker compose up --build
   ```

4. 打开 `http://localhost:8080`。后台登录账号来自 `.env` 中的 `ADMIN_USERNAME` 和 `ADMIN_PASSWORD`。

后台启动顺序由 Compose 控制：PostgreSQL 健康后启动 Backend，Backend 健康后再启动 Frontend。

### 后台运行与日志

```bash
docker compose up -d --build
docker compose ps
docker compose logs -f backend
```

### 停止服务

保留数据库数据：

```bash
docker compose down
```

删除服务及数据库数据卷：

```bash
docker compose down -v
```

`down -v` 会永久删除文章、评论和管理员数据，执行前请确认不再需要这些数据。

## 本地开发

### 后端

准备一个可访问的 PostgreSQL 数据库并设置必要环境变量，然后运行：

```bash
cd backend
go run ./cmd/server
```

后端默认监听 `8080`。启动时会自动执行 GORM 迁移，并在管理员表为空时初始化管理员。

### 前端

```bash
cd frontend
npm ci
npm run dev
```

开发地址默认为 `http://localhost:5173`。Vite 会将 `/api` 代理到 `http://127.0.0.1:8080`。如需直接指定 API 地址，可以设置：

```env
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

直接跨域访问后端时，请同时将前端 Origin 加入 `CORS_ALLOWED_ORIGINS`。

## 环境变量

| 变量 | 用途 | 示例说明 |
| --- | --- | --- |
| `APP_ENV` | 运行环境；生产环境使用 `production` | `development` |
| `POSTGRES_DB` | PostgreSQL 数据库名 | `personal_blog` |
| `POSTGRES_USER` | PostgreSQL 用户 | `blog_user` |
| `POSTGRES_PASSWORD` | PostgreSQL 密码 | 必须自行设置 |
| `POSTGRES_PORT` | 手工运行 PostgreSQL 时使用的端口提示 | `5432` |
| `SERVER_PORT` | 手工运行 Backend 时的监听端口 | `8080` |
| `DATABASE_URL` | GORM/PostgreSQL 连接字符串 | Docker 中主机名必须为 `postgres` |
| `ADMIN_USERNAME` | 首次初始化的管理员用户名 | `admin` |
| `ADMIN_PASSWORD` | 首次初始化的管理员密码，数据库只存 bcrypt 哈希 | 必须自行设置 |
| `JWT_SECRET` | JWT HMAC 签名密钥 | 至少 32 个字符 |
| `JWT_EXPIRES_HOURS` | JWT 有效时长（小时） | `24` |
| `CORS_ALLOWED_ORIGINS` | 允许跨域访问 API 的 Origin，逗号分隔且不带路径 | `https://blog.example.com` |
| `VITE_API_BASE_URL` | 可选的前端 API 基础地址 | 默认 `/api/v1` |

生产模式会拒绝 `.env.example` 中的示例管理员密码、数据库密码和 JWT Secret。

## 默认管理员初始化

Backend 每次启动都会检查 `administrators` 表：

- 表为空：使用 `ADMIN_USERNAME` 和 `ADMIN_PASSWORD` 创建唯一管理员，密码以 bcrypt 哈希保存。
- 已有管理员：不会覆盖账号或密码。

因此，在数据库已经初始化后仅修改 `.env` 中的管理员密码不会更新现有账号。若是一次尚无有效数据的本地安装，可以执行 `docker compose down -v` 后重新初始化；生产环境不要用清空数据卷的方式修改密码。

## API 简介

所有接口使用 `/api/v1` 前缀。

| 方法与路径 | 说明 | 鉴权 |
| --- | --- | --- |
| `GET /health` | 健康检查 | 否 |
| `GET /articles` | 已发布文章列表 | 否 |
| `GET /articles/:slug` | 已发布文章详情 | 否 |
| `GET /articles/:slug/comments` | 已审核评论列表 | 否 |
| `POST /articles/:slug/comments` | 提交待审核评论 | 否 |
| `POST /admin/auth/login` | 管理员登录 | 否 |
| `GET /admin/auth/me` | 当前管理员 | Bearer JWT |
| `GET /admin/dashboard` | 后台统计 | Bearer JWT |
| `/admin/articles...` | 文章管理 | Bearer JWT |
| `/admin/comments...` | 评论审核与删除 | Bearer JWT |

管理接口使用请求头：

```http
Authorization: Bearer <token>
```

## 测试与质量检查

后端：

```bash
cd backend
go test ./...
go vet ./...
```

前端：

```bash
cd frontend
npm run test
npm run typecheck
npm run build
```

检查 Compose 配置：

```bash
docker compose config --quiet
```

## Linux 服务器部署

1. 在服务器安装 Docker Engine 和 Docker Compose v2。
2. 上传项目或从 Git 仓库克隆。
3. 执行 `cp .env.example .env` 并设置生产变量：
   - `APP_ENV=production`
   - 使用强随机 `POSTGRES_PASSWORD`
   - 让 `DATABASE_URL` 使用相同数据库密码
   - 设置安全的 `ADMIN_PASSWORD`
   - 设置至少 32 字符的随机 `JWT_SECRET`
   - 将 `CORS_ALLOWED_ORIGINS` 设置为实际站点 Origin，例如 `https://blog.example.com`
4. 启动服务：

   ```bash
   docker compose up -d --build
   ```

5. 检查运行状态：

   ```bash
   docker compose ps
   curl http://localhost:8080/api/v1/health
   ```

PostgreSQL 和 Backend 只在 Compose 内部网络中通信；Compose 默认只向主机发布 Frontend 的 `8080` 端口。

### HTTPS

项目不自动申请证书。生产环境可以在 `8080` 端口外层使用以下任一方式终止 HTTPS：

- 服务器 Nginx
- Caddy
- 云厂商负载均衡或反向代理

外层代理应把请求完整转发到 `http://127.0.0.1:8080`，并保留 `Host`、`X-Forwarded-For` 和 `X-Forwarded-Proto` 等请求头。证书和 HTTPS 配置应放在服务器层，不要提交私钥或生产 `.env`。

## 常见错误排查

### `database connection failed`

- 检查 `POSTGRES_DB`、`POSTGRES_USER`、`POSTGRES_PASSWORD` 与 `DATABASE_URL` 是否一致。
- Docker 环境中的数据库主机名必须是 `postgres`，不能写 `localhost`。
- 如果修改了已有数据卷的 PostgreSQL 密码，旧数据卷不会自动采用新密码；请恢复原密码或在确认无数据后执行 `docker compose down -v`。

### 生产模式启动后立即退出

检查是否仍在使用 `.env.example` 中的 `change_me` 或示例 JWT Secret。生产模式会主动拒绝示例凭据。

### 管理员无法登录

- 确认使用的是数据库首次初始化时的账号密码。
- 修改 `.env` 不会覆盖已有管理员。
- 查看 `docker compose logs backend`，日志不会输出密码或完整 Token。

### 首页返回 `502 Bad Gateway`

运行 `docker compose ps` 确认 `postgres` 和 `backend` 已健康，再查看 `docker compose logs backend`。首次启动需要等待数据库健康检查和迁移完成。

### 刷新文章或后台页面出现 404

生产环境应使用仓库自带的 Frontend 镜像和 `deploy/nginx.conf`。如果使用其他 Web 服务器，需要把不存在的前端路径回退到 `index.html`，但不能对 `/api/` 做 History 回退。

### 跨域请求被浏览器拦截

将前端完整 Origin 加入 `CORS_ALLOWED_ORIGINS`，多个值用英文逗号分隔。Origin 只包含协议、域名和可选端口，不包含路径或末尾 `/`。

### `8080` 端口已被占用

修改 `docker-compose.yml` 中 Frontend 的端口映射，例如把 `8080:80` 改为 `8081:80`，然后访问新端口。
