package routers

import (
	"auth/controllers"

	"auth/lib/admin"

	"github.com/astaxie/beego"
)

func init() {
	admin.Run()
	beego.Router("/test", &controllers.MainController{})
	beego.Post("/login", controllers.Login)
	beego.Post("/logout", controllers.Logout)
	beego.Post("/chpwd", controllers.ChangePassword)
}
