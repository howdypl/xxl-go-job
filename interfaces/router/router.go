package router

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
)

// InitRouter 初始化路由，添加中间件
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	RegisterRouter(r)
	return r
}

// RegisterRouter 注册路由
func RegisterRouter(r *gin.Engine) {
	// 创建路由组
	g := r.Group("")
	g.GET("/", func(c *gin.Context) {
		c.JSON(
			http.StatusOK,
			gin.H{
				"code":  http.StatusOK,
				"message": "Welcome server",
			},
		)
	})

	// swagger；注意：生产环境可以注释掉
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

}
