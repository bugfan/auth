package main

import (
	_ "auth/routers"
	"runtime"

	"github.com/astaxie/beego"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
func main() {
	beego.Run()
}
