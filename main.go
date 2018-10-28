package main

import (
	"auth/db/redis"
	_ "auth/routers"
	"auth/util/jwt"
	"os"
	"runtime"

	"github.com/astaxie/beego"
	"github.com/bugfan/goini"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// 读环境变量
	if os.Getenv("ENV_FILE") != "" {
		goini.Env = goini.NewMyEnv(os.Getenv("ENV_FILE")) // 指定读取路径是否从指定文件读取
	} else {
		goini.Env = goini.NewMyEnv() // 不指定文件，读取默认环境变量
	}
}
func main() {
	redis.JWT = &redis.Redis{} // 链接redis
	redis.JWT.ConnRedis(("auth"))
	jwt.InitJWTConf() // 初始化配置
	beego.Run()
}
