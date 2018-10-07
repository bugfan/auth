package routers

import (
	"auth/controllers"

	"auth/lib/admin"

	"github.com/astaxie/beego"
)

func init() {
	admin.Run()
	beego.Router("/test", &controllers.MainController{})
	beego.Router("/api/jwt/user/login", &controllers.JWTController{})
	beego.Router("/api/jwt/user/logout", &controllers.JWTController{})
	beego.Router("/api/jwt/user/chpwd", &controllers.JWTController{})
}
