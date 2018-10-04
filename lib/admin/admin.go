package admin

import (
	"fmt"
	"mime"
	"os"

	. "auth/lib/admin/src/lib"
	"auth/lib/admin/src/models"

	"github.com/astaxie/beego"
)

const VERSION = "0.1.1"

func Run() {
	//初始化
	initialize()

	fmt.Println("Starting....")

	fmt.Println("Start ok")
}
func initialize() {
	mime.AddExtensionType(".css", "text/css")
	//判断初始化参数
	initArgs()

	models.Connect()

	router()
	beego.AddFuncMap("stringsToJson", StringsToJson)
}
func initArgs() {
	args := os.Args
	for _, v := range args {
		if v == "-syncdb" {
			models.Syncdb()
			os.Exit(0)
		}
	}
}
