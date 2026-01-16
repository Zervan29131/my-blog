package routes

import (
	"my-blog-backend/controllers"
	"my-blog-backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS())

	v1 := r.Group("/api/v1")
	{
		// 公开接口
		v1.POST("/register", controllers.Register)
		v1.POST("/login", controllers.Login)
		v1.GET("/posts", controllers.GetPostList)
		v1.GET("/posts/:id", controllers.GetPostDetail)

		// 新增: 获取标签列表 (通常写文章时需要选择标签，或者侧边栏显示，设为公开)
		v1.GET("/tags", controllers.GetTagList)

		// 权限接口
		authorized := v1.Group("/")
		authorized.Use(middleware.JWTAuthMiddleware())
		{
			authorized.POST("/posts", controllers.CreatePost)
			authorized.PUT("/posts/:id", controllers.UpdatePost)
			authorized.DELETE("/posts/:id", controllers.DeletePost)

			// 新增: 标签管理 (只有登录用户/管理员能创建和删除)
			authorized.POST("/tags", controllers.CreateTag)
			authorized.DELETE("/tags/:id", controllers.DeleteTag)
		}
	}
	return r
}
