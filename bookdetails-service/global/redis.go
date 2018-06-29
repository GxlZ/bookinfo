package global

import (
	"github.com/go-redis/redis"
)

func newRedisClient() *redis.Client {
	c := redis.NewClient(&redis.Options{
		Addr:     Conf.Redis.Addr,
		Password: Conf.Redis.Password,
		DB:       Conf.Redis.DB,
	})

	_, err := c.Ping().Result()
	if err != nil {
		Logger.Fatalln(err)
	}

	return c
}
