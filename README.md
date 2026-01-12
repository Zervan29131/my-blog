my-blog/
├── backend/                  # Go 后端项目
│   ├── config/               # 配置加载逻辑 (viper)
│   │   └── config.go
│   ├── controllers/          # 控制器层: 处理 HTTP 请求与参数解析
│   │   ├── auth_controller.go
│   │   └── post_controller.go
│   ├── database/             # 基础设施层: 数据库初始化 (GORM)
│   │   └── database.go
│   ├── middleware/           # 中间件: JWT认证, CORS, 日志
│   │   ├── jwt_auth.go
│   │   └── cors.go
│   ├── models/               # 数据模型层: 定义数据库 Struct
│   │   └── models.go
│   ├── routes/               # 路由层: 注册 URL 与控制器的绑定
│   │   └── routes.go
│   ├── services/             # 业务逻辑层: 复杂的业务逻辑与数据库交互
│   │   ├── auth_service.go
│   │   └── post_service.go
│   ├── utils/                # 工具库: JWT生成, 密码加密等
│   │   ├── jwt.go
│   │   └── password.go
│   ├── config.yaml           # 配置文件
│   ├── go.mod                # Go 依赖管理
│   └── main.go               # 程序入口
│
└── frontend/                 # Vue 3 前端项目 (建议由 Vite CLI 生成)
├── public/
├── src/
│   ├── api/              # API 接口统一管理 (Axios)
│   ├── assets/           # 静态资源 (图片, CSS)
│   ├── components/       # 公共组件
│   ├── router/           # 路由配置
│   ├── stores/           # 状态管理 (Pinia)
│   ├── views/            # 页面视图 (Page)
│   ├── App.vue
│   └── main.ts
├── index.html
└── package.json



frontend/
├── node_modules/       # (运行 npm install 后生成的依赖包文件夹)
├── public/             # 静态资源文件夹
│   └── vite.svg
├── src/
│   ├── api/            # API 接口定义
│   │   ├── auth.ts     # (已创建) 登录注册接口
│   │   └── post.ts     # (建议创建) 文章相关接口
│   ├── components/     # 公共组件 (目前为空)
│   ├── router/         # 路由配置
│   │   └── index.ts    # (已创建) 路由入口
│   ├── stores/         # Pinia 状态管理
│   │   └── user.ts     # (已创建) 用户状态
│   ├── utils/          # 工具函数
│   │   └── request.ts  # (已创建) Axios 封装
│   ├── views/          # 页面组件
│   │   └── Login.vue   # (已创建) 登录页
│   ├── App.vue         # (已创建) 根组件
│   ├── main.ts         # (已创建) Vue 入口文件
│   └── style.css       # (Vite 默认生成的样式文件，可选)
├── index.html          # (已创建) 浏览器访问入口 HTML
├── package.json        # (已创建) 项目依赖配置
├── tsconfig.json       # (Vite 模板生成) TypeScript 配置文件
└── vite.config.ts      # (已创建) Vite 构建工具配置




博客后端 API 测试流程

前提: 确保后端服务已启动 (http://localhost:8080)，且数据库连接正常。

1. 注册 (创建第一个用户)

这是为了初始化数据。

URL: POST http://localhost:8080/api/v1/register

Body (JSON):

{
    "username": "admin",
    "password": "password123"
}


预期结果: 200 OK, {"message": "注册成功"}

检查: 查看数据库 users 表，应该多了一条数据，且 password 字段是一串乱码（加密后的哈希）。

2. 登录 (获取凭证)

测试密码校验和 JWT 生成。

URL: POST http://localhost:8080/api/v1/login

Body (JSON):

{
    "username": "admin",
    "password": "password123"
}


预期结果: 200 OK, 返回类似如下数据：

{
    "message": "登录成功",
    "token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6ImFkbWluIiwiaXNzIjoibXktYmxvZy1iYWNrZW5kIiwiZXhwIjoxNzY4MzE3NzU0fQ.Go8kSlHkKXJWX0yg6efBQpx6Fok8C6qYo3D6euPI5Jg"

}}


关键动作: 复制这个 Token 字符串，接下来的操作都需要它。

3. 创建文章 (测试鉴权 + 写入)

测试 JWT 中间件是否生效。

URL: POST http://localhost:8080/api/v1/posts

Headers:

Key: Authorization

Value: Bearer <你刚才复制的Token>  (注意 Bearer 和 Token 之间有个空格)

Body (JSON):

{
    "title": "我的第一篇博客",
    "content": "# Hello World\n这是Markdown内容",
    "summary": "这是摘要",
    "category_id": 1, 
    "author_id": 1
}


(注: 如果 category_id 1 不存在可能会报错，你可以先去数据库随便插一条 category 数据，或者暂时无视外键约束)

预期结果: 201 Created, 返回创建的文章对象。

测试失败情况: 如果不加 Header，应该返回 401 Unauthorized。

4. 获取文章列表 (测试公开接口)

不需要 Token。

URL: GET http://localhost:8080/api/v1/posts?page=1&page_size=5

预期结果: 200 OK, data 数组里应该包含你刚才创建的文章。

5. 修改文章 (测试更新逻辑)

URL: PUT http://localhost:8080/api/v1/posts/1 (假设文章ID是1)

Headers: 加上 Authorization: Bearer ...

Body (JSON):

{
    "title": "修改后的标题"
}


预期结果: 200 OK, {"message": "更新成功"}。

6. 删除文章 (测试删除逻辑)

URL: DELETE http://localhost:8080/api/v1/posts/1

Headers: 加上 Authorization: Bearer ...

预期结果: 200 OK。再次调用列表接口，该文章应消失（或 deleted_at 字段有值）。