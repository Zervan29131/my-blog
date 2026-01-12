package routes

import (
	"my-blog-backend/controllers"
	"my-blog-backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS()) // 启用跨域

	v1 := r.Group("/api/v1")
	{
		// 公开接口
		v1.POST("/register", controllers.Register) // 新增注册接口
		v1.POST("/login", controllers.Login)
		v1.GET("/posts", controllers.GetPostList)
		v1.GET("/posts/:id", controllers.GetPostDetail)

		// 权限接口
		authorized := v1.Group("/")
		authorized.Use(middleware.JWTAuthMiddleware())
		{
			authorized.POST("/posts", controllers.CreatePost)
			authorized.PUT("/posts/:id", controllers.UpdatePost)
			authorized.DELETE("/posts/:id", controllers.DeletePost)
		}
	}
	return r
}
