package redis

import (
	"github.com/go-redis/redis/v7"
	"github.com/spf13/viper"
	"runtime"
)

const (
	PackageHashFileKey = "mirror:packageHashFile"
	ProviderKey        = "mirror:providers"
	Dist               = "mirror:dist"
	updateTimeKey         = "mirror:updateTIme"
)

var redisClient *redis.Client

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.addr"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.database"),
		PoolSize: 100 * runtime.NumCPU(),
	})
	_, err := redisClient.Ping().Result()
	if err != nil {
		panic(err)
	}
}
