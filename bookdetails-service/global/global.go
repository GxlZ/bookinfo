package global

import (
	"bookinfo/bookdetails-service/models"
	"github.com/openzipkin/zipkin-go/reporter"
	"os"
	"strings"
)

var Logger logger

var BOOK_DB *db

var zipkinReporter reporter.Reporter

var Redis *redisClient

func init() {

	for _, v := range os.Args {
		if strings.Contains(v, "-test.v=true") {
			return
		}
	}

	Logger = newLogger()

	loadConf()

	BOOK_DB = newBookDB()

	Redis = newRedisClient()

	models.Migrate(BOOK_DB.DB)

	zipkinReporter = NewZipkinReporter()
}
