package com

import (
	"auth/jwt/model"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/bugfan/logrus"
	"github.com/bugfan/to"
	_ "github.com/lib/pq"
)

func NewOrm() orm.Ormer {
	return orm.NewOrm()
}
func InitOrm(conf string) {
	time.Sleep(1e9)
	c, err := url.Parse(conf)
	if err != nil {
		logrus.Error("数据库连接参数错误:", err, conf)
	}
	logrus.Println(c.Host, c.User, c.User.String(), c.Scheme)
	pass, _ := c.User.Password()
	//orm.RegisterDriver("postgres", orm.DR_Postgres)
	orm.RegisterDataBase("default", "postgres", "user="+c.User.Username()+" password="+pass+" dbname="+c.Scheme+
		" host="+strings.Split(c.Host, ":")[0]+" port="+strings.Split(c.Host, ":")[1]+" sslmode=disable")
	model.Syncdb()
	orm.RunSyncdb("default", false, true)
	orm.Debug = to.Bool(os.Getenv("ORM_DEBUG"))
	orm.DebugLog = orm.NewLog(&LogrusWriter{})
	logrus.Info("数据库初始化完成! 您使用的是：", c.User.String())
}
