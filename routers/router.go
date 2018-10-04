package routers

import (
	"auth/controllers"

	"auth/lib/admin"

	"github.com/astaxie/beego"
)

func init() {
	admin.Run()
	beego.Router("/", &controllers.MainController{})
}
