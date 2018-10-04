package main

import (
	_ "auth/routers"

	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}
