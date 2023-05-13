package server

import "github.com/gin-gonic/gin"

var (
	routerEngine *gin.Engine
)

// CreateRouterEngine 创建路由引擎并添加路由
func CreateRouterEngine() *gin.Engine {

	routerEngine = gin.Default()
	routerEngine.POST("/api/go-cqhttp", MessagePreprocessing)

	return routerEngine
}
