// Author: xufei
// Date: 2019-10-21

package route

import "github.com/gin-gonic/gin"

type App struct {
	router *gin.Engine
}

func New() *App {
	app := &App{}
	app.initRouter()

	return &App{}
}

func (a *App) initRouter() {
	a.router = gin.Default()
	//	获取在线用户列表
}

func (a *App) Run() {
	a.router.Run(":8084")
}
