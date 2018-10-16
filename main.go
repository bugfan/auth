package main

import (
	"auth/db/redis"
	_ "auth/routers"

	"github.com/astaxie/beego"
	"github.com/bugfan/goini"
)

func main() {
	goini.Env = goini.NewMyEnv(".env")
	redis.JWT = &redis.Redis{}
	redis.JWT.ConnRedis(goini.Env.Getenv("PROJECT_NAME"))

	beego.Run()
}
