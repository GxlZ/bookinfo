package global

import (
	"bookinfo/bookcomments-service/models"
	"github.com/openzipkin/zipkin-go/reporter"
)

var Logger logger

var BOOK_DB *db

var zipkinReporter reporter.Reporter

func init() {

	loadConf()

	Logger = newLogger()

	BOOK_DB = newBookDB()

	models.Migrate(BOOK_DB.DB)

	zipkinReporter = NewZipkinReporter()
}
