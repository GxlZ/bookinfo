package global

import (
	goredis "github.com/go-redis/redis"
	"context"
	"time"
	"github.com/pquerna/ffjson/ffjson"
)

type redisClient struct {
	*goredis.Client
}

func (this *redisClient) WarpGet(ctx context.Context, key string) *goredis.StringCmd {
	span, _, err := zipkin(
		"redis",
		ctx,
		zipkinOption{
			OptionType: ZIPKIN_OPTION_TAG,
			zipkinTag:  ZipkinTag{"redis get key", key},
		},
	)


	v:= this.Client.Get(key)

	if err == nil {
		defer func(v *goredis.StringCmd) {
			span.Annotate(time.Now(), "out redis")
			span.Tag("redis get value",v.Val())
			span.Finish()
		}(v)
	}

	return v

}

func (this *redisClient) WarpSet(ctx context.Context, key string, value interface{}, expiration time.Duration) *goredis.StatusCmd {
	bytes, _ := ffjson.Marshal(value)

	span, _, err := zipkin(
		"redis",
		ctx,
		zipkinOption{
			OptionType: ZIPKIN_OPTION_TAG,
			zipkinTag:  ZipkinTag{"redis set key: " + key, string(bytes)},
		},
	)
	if err == nil {
		defer func() {
			span.Annotate(time.Now(), "out redis")
			span.Finish()
		}()
	}

	return this.Client.Set(key, bytes, expiration)
}

func newRedisClient() *redisClient {
	rc := goredis.NewClient(&goredis.Options{
		Addr:     Conf.Redis.Addr,
		Password: Conf.Redis.Password,
		DB:       Conf.Redis.DB,
	})

	_, err := rc.Ping().Result()
	if err != nil {
		Logger.Fatalln(err)
	}

	return &redisClient{
		rc,
	}
}
