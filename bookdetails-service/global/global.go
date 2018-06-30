package global

import (
	"bookinfo/bookdetails-service/models"
	"github.com/openzipkin/zipkin-go/reporter"
)

var Logger logger

var BOOK_DB *db

var zipkinReporter reporter.Reporter

var Redis *redisClient

func init() {

	Logger = newLogger()

	loadConf()

	BOOK_DB = newBookDB()

	Redis = newRedisClient()

	models.Migrate(BOOK_DB.DB)

	zipkinReporter = NewZipkinReporter()
}
