# 个人博客系统搭建需求

## 1. 项目目标

开发一个前后端分离的个人博客系统，满足以下核心需求：

1. 访客可以浏览已发布文章。
2. 访客可以在文章下提交评论。
3. 管理员可以登录后台管理文章和评论。
4. 项目当前可在本地运行，并可通过 Docker Compose 平滑部署到 Linux 服务器。
5. 优先保证结构简单、可维护、可运行，不做过度设计。

本项目定位为个人使用的 MVP，不按大型内容平台设计。

---

## 2. 固定技术栈

除非出现无法解决的兼容问题，否则不得更换以下技术栈。

### 2.1 前端

* Vue 3
* TypeScript
* Vite
* Vue Router
* Pinia
* Axios
* Element Plus
* markdown-it
* DOMPurify

### 2.2 后端

* Go
* Gin
* GORM
* PostgreSQL
* JWT
* bcrypt
* Go 标准库 `log/slog`

### 2.3 部署

* Docker
* Docker Compose
* Nginx
* PostgreSQL 持久化数据卷

### 2.4 接口规范

* REST API
* JSON 数据格式
* API 统一使用 `/api/v1` 前缀
* UTF-8 编码
* 时间统一由后端使用 ISO 8601 格式返回

---

## 3. MVP 功能范围

### 3.1 博客前台

#### 首页

展示已发布文章列表。

每篇文章显示：

* 标题
* 摘要
* 发布时间
* 评论数量
* 查看全文入口

支持分页：

* 默认每页 10 篇
* 最大每页 50 篇

#### 文章详情页

展示：

* 文章标题
* 发布时间
* Markdown 渲染后的正文
* 已审核通过的评论
* 评论提交表单

文章正文在数据库中保存为 Markdown 原文，在前端渲染为 HTML。

所有渲染后的 HTML 必须通过 DOMPurify 清理，避免 XSS。

#### 评论提交

访客不需要注册账号。

评论字段：

* 昵称，必填，2～50 个字符
* 邮箱，选填，最大 255 个字符
* 评论内容，必填，2～1000 个字符

评论提交后默认状态为 `pending`，只有管理员审核通过后才在前台显示。

提交成功后提示：

> 评论已提交，审核通过后将会显示。

邮箱不得在博客前台展示，也不得通过公开接口返回。

---

### 3.2 管理后台

#### 管理员登录

系统仅支持一个管理员账号。

不实现：

* 用户注册
* 找回密码
* 多管理员
* 角色权限
* 第三方登录
* Refresh Token

管理员账号通过环境变量初始化：

```env
ADMIN_USERNAME=admin
ADMIN_PASSWORD=请设置安全密码
JWT_SECRET=至少32位随机字符串
```

首次启动时，如果管理员表为空，则创建管理员账号。

数据库中只能保存 bcrypt 密码哈希，不得保存明文密码。

JWT 有效期为 24 小时。JWT 过期后重新登录。

#### 后台首页

显示简单统计数据：

* 文章总数
* 已发布文章数
* 草稿文章数
* 待审核评论数
* 已通过评论数

不需要图表。

#### 文章管理

管理员可以：

* 查看全部文章
* 创建文章
* 编辑文章
* 删除文章
* 保存为草稿
* 发布文章
* 将已发布文章改回草稿

文章编辑字段：

* 标题
* Slug
* 摘要
* Markdown 正文
* 状态

编辑器采用普通 Markdown 文本区域和预览区域，不集成复杂富文本编辑器。

Slug 规则：

* 只能包含小写英文字母、数字和连字符
* 必须唯一
* 创建文章时可以由标题自动生成
* 管理员可以手动修改
* Slug 冲突时接口返回明确错误

删除文章时，同时删除该文章下的评论。

删除操作必须在前端弹出确认框。

#### 评论管理

管理员可以：

* 查看全部评论
* 按状态筛选评论
* 查看评论所属文章
* 审核通过评论
* 拒绝评论
* 删除评论

评论状态：

* `pending`
* `approved`
* `rejected`

MVP 不实现评论回复和嵌套评论。

---

## 4. 明确不包含的功能

为了控制开发范围和 Codex 消耗，以下功能不在本期实现：

* 用户注册和用户中心
* 多管理员和权限系统
* 分类和标签
* 站内全文搜索
* 文章图片上传
* 对象存储
* 嵌套评论
* 评论回复
* 邮件通知
* 验证码
* OAuth
* RSS
* 点赞和收藏
* 浏览量统计
* 国际化
* SSR
* Redis
* 消息队列
* 微服务
* Kubernetes
* 自动生成 SEO sitemap
* 主题切换
* 自定义站点设置后台

不得自行增加上述功能。

---

## 5. 数据库设计

### 5.1 administrators 表

| 字段            | 类型           | 约束    |
| ------------- | ------------ | ----- |
| id            | bigint       | 主键，自增 |
| username      | varchar(50)  | 唯一，非空 |
| password_hash | varchar(255) | 非空    |
| created_at    | timestamp    | 非空    |
| updated_at    | timestamp    | 非空    |

### 5.2 articles 表

| 字段           | 类型           | 约束                |
| ------------ | ------------ | ----------------- |
| id           | bigint       | 主键，自增             |
| title        | varchar(200) | 非空                |
| slug         | varchar(200) | 唯一，非空             |
| summary      | varchar(500) | 非空，允许空字符串         |
| content      | text         | 非空                |
| status       | varchar(20)  | draft 或 published |
| published_at | timestamp    | 可空                |
| created_at   | timestamp    | 非空                |
| updated_at   | timestamp    | 非空                |

业务规则：

* 新文章默认状态为 `draft`。
* 第一次发布时设置 `published_at`。
* 修改已发布文章时保留原发布时间。
* 草稿文章不得通过公开接口访问。
* 访问不存在或未发布的文章时返回 404。

### 5.3 comments 表

| 字段         | 类型            | 约束                          |
| ---------- | ------------- | --------------------------- |
| id         | bigint        | 主键，自增                       |
| article_id | bigint        | 外键，非空                       |
| nickname   | varchar(50)   | 非空                          |
| email      | varchar(255)  | 可空                          |
| content    | varchar(1000) | 非空                          |
| status     | varchar(20)   | pending、approved 或 rejected |
| created_at | timestamp     | 非空                          |
| updated_at | timestamp     | 非空                          |

业务规则：

* 删除文章时级联删除评论。
* 公开接口只返回 `approved` 评论。
* 公开接口不得返回邮箱字段。

---

## 6. API 设计

### 6.1 统一响应格式

成功响应：

```json
{
  "data": {},
  "message": "success"
}
```

分页响应：

```json
{
  "data": {
    "items": [],
    "page": 1,
    "page_size": 10,
    "total": 0,
    "total_pages": 0
  },
  "message": "success"
}
```

错误响应：

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "请求参数不正确",
    "details": {}
  }
}
```

不得在生产响应中返回：

* SQL 语句
* 数据库错误详情
* Go 调用栈
* JWT 密钥
* 环境变量
* 内部文件路径

---

### 6.2 公共接口

#### 健康检查

```http
GET /api/v1/health
```

响应：

```json
{
  "data": {
    "status": "ok"
  },
  "message": "success"
}
```

#### 获取已发布文章列表

```http
GET /api/v1/articles?page=1&page_size=10
```

返回字段：

* id
* title
* slug
* summary
* published_at
* comment_count

列表接口不得返回完整正文。

#### 获取文章详情

```http
GET /api/v1/articles/:slug
```

返回字段：

* id
* title
* slug
* summary
* content
* published_at
* created_at
* updated_at

#### 获取文章评论

```http
GET /api/v1/articles/:slug/comments?page=1&page_size=20
```

仅返回审核通过的评论：

* id
* nickname
* content
* created_at

不得返回邮箱。

#### 提交评论

```http
POST /api/v1/articles/:slug/comments
Content-Type: application/json
```

请求：

```json
{
  "nickname": "访客",
  "email": "visitor@example.com",
  "content": "评论内容"
}
```

新评论状态固定为 `pending`，前端不得传入状态。

---

### 6.3 管理员认证接口

#### 登录

```http
POST /api/v1/admin/auth/login
```

请求：

```json
{
  "username": "admin",
  "password": "password"
}
```

响应：

```json
{
  "data": {
    "token": "jwt-token",
    "expires_in": 86400
  },
  "message": "success"
}
```

#### 获取当前管理员

```http
GET /api/v1/admin/auth/me
Authorization: Bearer <token>
```

除登录接口外，所有 `/admin` 接口都必须验证 JWT。

---

### 6.4 后台统计接口

```http
GET /api/v1/admin/dashboard
Authorization: Bearer <token>
```

响应字段：

* article_total
* article_published
* article_draft
* comment_pending
* comment_approved

---

### 6.5 后台文章接口

```http
GET    /api/v1/admin/articles
GET    /api/v1/admin/articles/:id
POST   /api/v1/admin/articles
PUT    /api/v1/admin/articles/:id
DELETE /api/v1/admin/articles/:id
```

文章创建请求：

```json
{
  "title": "文章标题",
  "slug": "article-slug",
  "summary": "文章摘要",
  "content": "# Markdown 内容",
  "status": "draft"
}
```

文章更新使用完整更新方式，不同时实现 PUT 和 PATCH。

---

### 6.6 后台评论接口

```http
GET    /api/v1/admin/comments?status=pending&page=1&page_size=20
PUT    /api/v1/admin/comments/:id/status
DELETE /api/v1/admin/comments/:id
```

更新状态请求：

```json
{
  "status": "approved"
}
```

只允许：

* pending
* approved
* rejected

---

## 7. HTTP 状态码

统一使用以下状态码：

* `200`：查询或更新成功
* `201`：创建成功
* `204`：删除成功
* `400`：参数错误
* `401`：未登录或 Token 无效
* `404`：资源不存在
* `409`：Slug 或用户名冲突
* `429`：请求频率过高
* `500`：服务器内部错误

---

## 8. 前端路由

### 8.1 公共页面

```text
/                       文章列表
/articles/:slug         文章详情和评论
```

### 8.2 管理后台

```text
/admin/login             管理员登录
/admin                   后台统计
/admin/articles          文章列表
/admin/articles/new      新建文章
/admin/articles/:id/edit 编辑文章
/admin/comments          评论管理
```

后台路由需要前端路由守卫。

没有 Token 时访问后台页面，跳转到 `/admin/login`。

后端返回 401 时：

1. 删除本地 Token。
2. 跳转到登录页。
3. 不无限重试请求。

---

## 9. 前端状态管理

Pinia 只保存必要状态：

### authStore

* token
* currentAdmin
* isAuthenticated
* login()
* logout()
* fetchCurrentAdmin()

Token MVP 阶段保存在 `localStorage`。

不要将文章列表、评论列表等普通请求数据全部放进 Pinia，页面组件直接请求即可。

---

## 10. 后端目录结构

保持简单，不使用复杂的领域驱动设计。

```text
backend/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   ├── database/
│   ├── handler/
│   ├── middleware/
│   ├── model/
│   ├── router/
│   └── service/
├── go.mod
├── go.sum
└── Dockerfile
```

约束：

* Handler 负责解析请求和返回响应。
* Service 负责主要业务逻辑。
* Model 负责数据库模型。
* 不创建 Repository 接口层。
* 不引入依赖注入框架。
* 不使用代码生成框架。
* 不拆分为多个 Go 服务。

---

## 11. 前端目录结构

```text
frontend/
├── src/
│   ├── api/
│   ├── components/
│   ├── layouts/
│   ├── router/
│   ├── stores/
│   ├── types/
│   ├── utils/
│   ├── views/
│   │   ├── public/
│   │   └── admin/
│   ├── App.vue
│   └── main.ts
├── index.html
├── package.json
├── tsconfig.json
├── vite.config.ts
└── Dockerfile
```

要求：

* API 地址从环境变量读取。
* Axios 使用统一实例。
* 统一处理 401 和通用错误。
* TypeScript 不使用大量 `any`。
* 页面组件不直接拼接重复的 API 基础地址。
* 后台页面使用 Element Plus。
* 前台样式保持简洁，不建设复杂设计系统。

---

## 12. 项目根目录结构

```text
personal-blog/
├── backend/
├── frontend/
├── deploy/
│   └── nginx.conf
├── docs/
│   └── REQUIREMENTS.md
├── .env.example
├── .gitignore
├── docker-compose.yml
└── README.md
```

---

## 13. 环境变量

根目录提供 `.env.example`：

```env
APP_ENV=development

POSTGRES_DB=personal_blog
POSTGRES_USER=blog_user
POSTGRES_PASSWORD=change_me
POSTGRES_PORT=5432

SERVER_PORT=8080
DATABASE_URL=host=postgres user=blog_user password=change_me dbname=personal_blog port=5432 sslmode=disable

ADMIN_USERNAME=admin
ADMIN_PASSWORD=change_me
JWT_SECRET=replace_with_a_random_string_at_least_32_characters

JWT_EXPIRES_HOURS=24
CORS_ALLOWED_ORIGINS=http://localhost:5173,http://localhost:8080
```

要求：

* `.env` 不得提交到 Git。
* 仓库中不得出现真实密码。
* 应用启动时检查必要环境变量。
* 缺少关键变量时输出清晰错误并退出。
* 生产环境不得使用示例密码。

---

## 14. Docker 与部署要求

`docker-compose.yml` 包含三个服务：

```text
postgres
backend
frontend
```

其中：

* PostgreSQL 使用持久化 volume。
* Backend 等待数据库健康后启动。
* Frontend 构建后由 Nginx 提供静态文件。
* Nginx 将 `/api/` 转发到 Backend。
* Vue 路由刷新时统一回退到 `index.html`。
* 对外只需暴露一个 Web 端口。
* PostgreSQL 默认不暴露到公网。

本地启动命令：

```bash
docker compose up --build
```

预期访问地址：

```text
http://localhost:8080
```

停止服务：

```bash
docker compose down
```

停止并清除数据库：

```bash
docker compose down -v
```

迁移到服务器时只需要：

1. 安装 Docker 和 Docker Compose。
2. 上传或克隆项目。
3. 创建生产 `.env`。
4. 执行 `docker compose up -d --build`。
5. 后续可在服务器外层配置 HTTPS。

MVP 不要求自动申请 HTTPS 证书，但 README 需要说明可以使用 Nginx、Caddy 或云厂商反向代理配置 HTTPS。

---

## 15. 安全要求

必须实现：

* bcrypt 密码哈希
* JWT 签名验证
* 后台接口鉴权
* 请求参数长度验证
* Markdown HTML 清理
* GORM 参数化查询
* 评论接口基础限流
* CORS 白名单
* 敏感字段不出现在公开响应中
* 生产错误不返回内部堆栈

评论接口限流：

* 同一 IP 每分钟最多提交 5 次
* 使用进程内存实现即可
* 不引入 Redis
* 服务重启后限流状态清空可以接受

暂不实现验证码。

---

## 16. 日志要求

使用 Go 标准库 `log/slog`。

日志至少包含：

* 服务启动
* 数据库连接结果
* 请求方法
* 请求路径
* HTTP 状态码
* 请求耗时
* 管理员登录失败
* 服务器内部错误

禁止记录：

* 管理员密码
* JWT 完整内容
* 评论者邮箱
* 数据库密码
* JWT Secret

---

## 17. 测试要求

### 17.1 后端

至少覆盖：

1. 管理员正确密码登录成功。
2. 管理员错误密码登录失败。
3. 未携带 Token 访问后台接口返回 401。
4. 创建草稿文章成功。
5. 草稿文章无法通过公开接口访问。
6. 发布文章后可以通过公开接口访问。
7. 重复 Slug 返回 409。
8. 评论提交后状态为 pending。
9. 未审核评论不出现在公开评论列表。
10. 审核通过后评论可以公开查看。
11. 公开评论响应不包含 email。

### 17.2 前端

不要求完整单元测试体系，但必须通过：

```bash
npm run typecheck
npm run build
```

### 17.3 后端质量检查

必须通过：

```bash
go test ./...
go vet ./...
```

---

## 18. README 要求

README 必须包含：

* 项目简介
* 技术栈
* 目录结构
* 环境要求
* 本地启动方式
* Docker 启动方式
* 环境变量说明
* 默认管理员初始化逻辑
* 数据清理方式
* 服务器部署步骤
* API 简介
* 常见错误排查

README 中不得写入真实密码。

---

## 19. 验收标准

项目完成后必须满足以下条件：

1. 执行 `docker compose up --build` 可以启动全部服务。
2. 浏览器访问 `http://localhost:8080` 可以打开博客首页。
3. 管理员可以使用环境变量中的账号密码登录。
4. 管理员可以创建一篇草稿。
5. 草稿不会出现在博客首页。
6. 管理员发布文章后，文章出现在首页。
7. 访客可以查看文章 Markdown 正文。
8. 访客可以提交评论。
9. 新评论不会立即公开显示。
10. 管理员可以审核评论。
11. 审核通过后评论出现在文章详情页。
12. 管理员可以删除评论。
13. 管理员可以编辑和删除文章。
14. 刷新任意 Vue 页面不会出现 Nginx 404。
15. PostgreSQL 容器重启后数据仍然存在。
16. 后端测试、Go Vet、前端类型检查和前端构建全部通过。
17. 项目可以在另一台安装了 Docker 的 Linux 服务器上运行。
18. 仓库不包含密码、Token 或真实 `.env` 文件。

---

## 20. Codex 实施顺序

Codex 必须严格按以下任务顺序实施，每次只执行一个任务。

### TASK-01：项目骨架

完成：

* 根目录结构
* Vue 项目初始化
* Go 项目初始化
* PostgreSQL Docker 服务
* Backend Dockerfile
* Frontend Dockerfile
* Nginx 配置
* Docker Compose
* 健康检查接口
* `.env.example`
* `.gitignore`

验收：

```bash
docker compose up --build
curl http://localhost:8080/api/v1/health
```

---

### TASK-02：数据库和管理员认证

完成：

* 数据库连接
* GORM 模型
* 自动建表
* 管理员初始化
* bcrypt 密码校验
* JWT 生成和验证
* 登录接口
* 当前管理员接口
* 鉴权中间件
* 认证相关测试

不要实现文章和评论前端。

---

### TASK-03：文章后端接口

完成：

* 文章模型
* 后台文章 CRUD
* 发布和草稿逻辑
* Slug 校验
* Slug 唯一性处理
* 公开文章列表
* 公开文章详情
* 分页
* 文章相关测试

---

### TASK-04：评论后端接口

完成：

* 评论模型
* 公开评论提交
* 公开评论列表
* 后台评论列表
* 评论状态更新
* 评论删除
* 评论限流
* 邮箱字段保护
* 评论相关测试

---

### TASK-05：博客前台

完成：

* 首页文章列表
* 分页
* 文章详情
* Markdown 渲染
* DOMPurify 清理
* 评论列表
* 评论提交表单
* 基础加载状态
* 基础错误提示
* 移动端可用布局

不要提前实现额外功能。

---

### TASK-06：管理后台

完成：

* 登录页面
* Token 管理
* 路由守卫
* 后台布局
* 统计页面
* 文章列表
* 新建文章
* 编辑文章
* Markdown 预览
* 删除文章确认
* 评论列表
* 评论筛选
* 评论审核
* 评论删除

---

### TASK-07：联调与交付

完成：

* 修复接口联调问题
* 检查 Docker 网络
* 检查 Nginx 路由
* 检查 Vue History 回退
* 执行全部测试和构建
* 完成 README
* 检查敏感信息
* 验证全新环境启动
* 不增加新功能

---

## 21. Codex 工作约束

Codex 在每个任务中必须遵守：

1. 开始前读取本文件。
2. 先检查已有代码，不重复创建已有实现。
3. 一次只完成指定任务。
4. 不实现后续任务。
5. 不改变固定技术栈。
6. 不增加非需求功能。
7. 不引入微服务、Redis、消息队列或复杂架构。
8. 优先修改现有文件，不进行无关重构。
9. 依赖必须有明确用途。
10. 不生成大量无意义注释。
11. 不为简单逻辑创建多层抽象。
12. 不伪造测试通过结果。
13. 完成后实际执行相关测试或构建命令。
14. 遇到错误先定位并修复，不通过删除功能规避。
15. 最终只报告：

    * 修改了哪些文件
    * 完成了哪些内容
    * 执行了哪些命令
    * 测试是否通过
    * 是否存在阻塞问题

---

## 22. 首次交给 Codex 的提示词

```text
请先完整阅读 docs/REQUIREMENTS.md。

这是项目的唯一需求来源。不要自行扩展功能，不要更换技术栈，不要提前实现后续任务。

现在只执行 TASK-01：项目骨架。

执行要求：
1. 先检查当前仓库已有文件。
2. 按 REQUIREMENTS.md 创建最小可运行结构。
3. 不实现文章、评论和完整认证功能。
4. 完成后实际运行 Docker 构建和健康检查。
5. 修复本任务范围内的错误。
6. 最后简要列出修改文件、验证命令和验证结果。
```

---

## 23. 后续任务提示词模板

后续每次只需要向 Codex 提供以下短提示，不要重复粘贴全部需求：

```text
请读取 docs/REQUIREMENTS.md，并检查当前仓库已有实现。

现在只执行 TASK-XX。

严格遵守该任务边界，不实现后续任务，不增加需求外功能，不更换技术栈。

完成后运行该任务要求的测试或构建命令，并只报告：
1. 修改文件
2. 完成内容
3. 验证命令
4. 验证结果
5. 阻塞问题
```

将 `TASK-XX` 替换为当前任务编号，例如：

```text
现在只执行 TASK-03。
```

---

## 24. 修复问题时的低额度提示词

发现 Bug 时，不要重新要求 Codex 审查整个项目，使用以下格式：

```text
请读取 docs/REQUIREMENTS.md。

当前问题：
<粘贴错误现象或错误日志>

预期行为：
<描述预期结果>

限制：
1. 只修复该问题。
2. 不重构无关代码。
3. 不增加新功能。
4. 先定位根因，再实施最小修改。
5. 修复后运行与该问题直接相关的测试或构建命令。
6. 最后只说明根因、修改文件和验证结果。
```

---

## 25. 项目完成定义

只有同时满足以下条件，项目才视为完成：

* TASK-01 至 TASK-07 全部完成
* 所有验收标准通过
* Docker Compose 可从空环境启动
* 数据持久化正常
* 前后台核心流程可用
* 测试和构建通过
* README 完整
* 不存在真实密钥
* 没有实现需求范围之外的大型功能
