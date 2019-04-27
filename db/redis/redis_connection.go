package redis

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/bugfan/goini"

	"gopkg.in/redis.v5"
)

var JWT *Redis // web jwt 链接

func init() {
	JWT = new(Redis)
}

// Redis redis连接对象
type Redis struct {
	RedisCluster *redis.Client
}

// ConnRedis 获取redis集群的链接
func (rc *Redis) ConnRedis(project string) {
	project = strings.ToUpper(project) + "_"
	addrs := rc.getRedisAddrs(project)
	if len(addrs) < 1 {
		log.Fatal("Redis No Host! Please Check ENV")
	}
	p := goini.Env.Getenv(project + "REDIS_PASSWORD")
	poolSize, _ := strconv.Atoi(goini.Env.Getenv(project + "REDIS_POOL_SIZE"))
	ind, _ := strconv.Atoi(goini.Env.Getenv(project + "REDIS_INDEX"))
	options := redis.Options{
		Addr: addrs[0],
		// MaxRedirects: 16, //最大重试次数，默认16
		PoolSize:    poolSize,
		DB:          ind,
		ReadTimeout: 500 * time.Millisecond,
		IdleTimeout: 12 * time.Second,
	}
	if strings.TrimSpace(p) != "" {
		options.Password = p
	}
	rc.RedisCluster = redis.NewClient(&options)
	if rc.Set("bar", 1, 1000000000) {
		log.Println("Redis Init Succeeded!")
	} else {
		log.Fatal("Redis Init Failed!")
	}
}

func (rc *Redis) getRedisAddrs(project string) []string {
	env, _ := rc.readEnv()
	var addrs []string
	for k, v := range env {
		if strings.HasPrefix(k, project+"REDIS_HOST") {
			addrs = append(addrs, v)
		}
	}
	return addrs
}
func (rc *Redis) readEnv() (map[string]string, error) {
	m := goini.Env.GetAll()
	return m, nil
}
