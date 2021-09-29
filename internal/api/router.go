/**
 * @File: api.go
 * @Author: hsien
 * @Description:
 * @Date: 9/18/21 3:01 PM
 */

package api

import (
	"custom_server/internal/api/middleware"
	"custom_server/internal/api/user"
	"custom_server/pkg/server"
	"custom_server/pkg/util"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

func ServerRouter(app *server.GinApp) {
	engine := app.Engine()
	engine.Use(gin.Recovery(),
		requestid.New(requestid.Config{
			Generator: func() string {
				return util.GenerateRandStr(8)
			}}))
	if app.DebugMode() {
		engine.Use(middleware.LoggerWithBody())
	} else {
		engine.Use(middleware.Logger())
	}

	// cron demo
	//print := cronjob.NewPrint()
	//app.AddCronJob(print)

	// user router
	userGroup := engine.Group("/user")
	user := user.NewUser(app.DB("postgres"))
	userGroup.GET("/list", middleware.Auth(), Wrap(user.UserList))
	userGroup.GET("/info", Wrap(user.UserInfo))
	userGroup.POST("/login", Wrap(user.Login))
	userGroup.POST("/register", Wrap(user.Register))
}
