package global

import (
	"bookinfo/bookdetails-service/models"
	"github.com/openzipkin/zipkin-go/reporter"
	"github.com/go-redis/redis"
)

var Logger logger

var BOOK_DB *db

var zipkinReporter reporter.Reporter

var Redis *redis.Client

func init() {

	Logger = newLogger()

	loadConf()

	BOOK_DB = newBookDB()

	Redis = newRedisClient()

	models.Migrate(BOOK_DB.DB)

	zipkinReporter = NewZipkinReporter()
}
