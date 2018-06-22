package global

import (
	"bookinfo/bookdetails-service/models"
	"github.com/openzipkin/zipkin-go/reporter"
)

var Logger logger

var BOOK_DB *db

var zipkinReporter reporter.Reporter

func init() {

	Logger = newLogger()

	loadConf()

	BOOK_DB = newBookDB()

	models.Migrate(BOOK_DB.DB)

	zipkinReporter = NewZipkinReporter()
}
