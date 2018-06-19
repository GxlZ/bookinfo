package global

import (
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"context"
	"time"
	zipkingo "github.com/openzipkin/zipkin-go"
	zipkingomodel "github.com/openzipkin/zipkin-go/model"
)

type db struct {
	*gorm.DB
}

const (
	ZIPKIN_OPTION_TAG      = 1 + iota
	ZIPKIN_OPTION_ANNOTATE
)

type zipkinOption struct {
	OptionType int
	zipkinTag  ZipkinTag
	Annotate   zipkingomodel.Annotation
}
type ZipkinTag struct {
	K string
	V string
}

func (this *db) zipkin(ctx context.Context, opts ...zipkinOption) (zipkingo.Span, context.Context, error) {
	span, newCtx, err := NewZipkinSpanFromCtx(ctx, func() (*zipkingo.Tracer, error) {
		return NewZipkinTracer()
	})

	if err != nil {
		return span, newCtx, err
	}

	span.Annotate(time.Now(), "in db")

	for _, option := range opts {
		switch option.OptionType {
		case ZIPKIN_OPTION_TAG:
			span.Tag(option.zipkinTag.K, option.zipkinTag.V)
		case ZIPKIN_OPTION_ANNOTATE:
			span.Annotate(option.Annotate.Timestamp, option.Annotate.Value)
		}
	}

	span.SetName("execute sql")

	return span, newCtx, err
}

func (this *db) WarpFind(ctx context.Context, out interface{}, where ...interface{}) (*gorm.DB) {
	span, _, err := this.zipkin(ctx)
	if err == nil {
		defer func() {
			span.Annotate(time.Now(), "out db")
			span.Finish()
		}()
	}

	return this.DB.Find(out, where)
}

func (this *db) WarpRawScan(ctx context.Context, dest interface{}, sql string, values ...interface{}) (*gorm.DB) {
	span, _, err := this.zipkin(
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
