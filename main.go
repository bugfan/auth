package main

import (
	"auth/db/redis"
	_ "auth/routers"
	"auth/util/jwt"
	"log"
	"os"
	"runtime"

	"github.com/astaxie/beego"
	"github.com/bugfan/goini"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// ready env
	if os.Getenv("AUTH_ENV") != "" {
		goini.Env = goini.NewEnv(os.Getenv("AUTH_ENV")) // read from file
	} else {
		goini.Env = goini.NewEnv() // read from sys env
	}
}
func main() {
	redis.JWT.ConnRedis("auth") // link auth redis
	jwt.InitJWTConf()           // init jwt config
	j, ee := jwt.GetJWT3(map[string]string{"zxy": "90"})
	log.Println("JWT:", ee, j)
	beego.Run()
}
