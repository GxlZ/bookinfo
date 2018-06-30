package global

import (
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"context"
	"time"
)

type db struct {
	*gorm.DB
}

func (this *db) WarpFind(ctx context.Context, out interface{}, where ...interface{}) (*gorm.DB) {
	span, _, err := zipkin("db",ctx)
	if err == nil {
		defer func() {
			span.Annotate(time.Now(), "out db")
			span.Finish()
		}()
	}

	return this.DB.Find(out, where)
}

func (this *db) WarpRawScan(ctx context.Context, dest interface{}, sql string, values ...interface{}) (*gorm.DB) {
	span, _, err := zipkin(
		"db",
		ctx,
		zipkinOption{
			OptionType: ZIPKIN_OPTION_TAG,
			zipkinTag:  ZipkinTag{"sql", sql},
		},
		zipkinOption{
			OptionType: ZIPKIN_OPTION_TAG,
			zipkinTag:  ZipkinTag{"values", fmt.Sprint(values)},
		},
	)
	if err == nil {
		defer func() {
			span.Annotate(time.Now(), "out db")
			span.Finish()
		}()
	}

	return this.DB.Raw(sql, values).Scan(dest)
}

func newBookDB() *db {
	conf := Conf.DB_BOOK
	entity, err := gorm.Open(
		conf.Driver,
		fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%s&loc=%s",
			conf.Username,
			conf.Password,
			conf.Host,
			conf.Port,
			conf.DBName,
			conf.Charset,
			conf.ParseTime,
			conf.Local,
		),
	)

	if err != nil {
		Logger.Fatalln("mysql conn failed,", err)
	}

	return &db{
		entity,
	}
}
