package global

import (
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
	"fmt"
)

var logFormatter = &logrus.TextFormatter{
	DisableColors:    false,
	FullTimestamp:    true,
	TimestampFormat:  "2006-01-02 15:04:05.0000",
	QuoteEmptyFields: true,
}

type logger struct {
	*logrus.Logger
}

func (this logger) Console(fieldsFunc FieldsFunc, args ...interface{}) {
	l := logrus.New()
	l.Formatter = logFormatter
	entry := l.WithFields(fieldsFunc().Get())
	entry.Data["file"] = fileInfo(2)
	entry.Println(args)
}

func (this logger) Log(keyvals ...interface{}) error {
	this.LogWithFields(func() *LogFields {
		return NewLogFields()
	}, keyvals)
	return nil
}

func (this logger) LogWithFields(fieldsFunc FieldsFunc, keyvals ...interface{}) error {
	entry := this.Logger.WithFields(fieldsFunc().Get())
	entry.Data["file"] = fileInfo(3)
	entry.Info(keyvals)
	return nil
}

func (this logger) Info(args ...interface{}) {
	this.InfoWithFields(func() *LogFields {
		return NewLogFields()
	}, args)
}

func (this logger) InfoWithFields(fieldsFunc FieldsFunc, args ...interface{}) {
	entry := this.WithFields(fieldsFunc().Get())
	entry.Data["file"] = fileInfo(3)
	entry.Infoln(args...)
}

func (this logger) Errorln(args ...interface{}) {
	this.ErrorlnWithFields(func() *LogFields {
		return NewLogFields()
	}, args)
}

func (this logger) ErrorlnWithFields(fieldsFunc FieldsFunc, args ...interface{}) {
	entry := this.WithFields(fieldsFunc().Get())
	entry.Data["file"] = fileInfo(3)
	entry.Errorln(args...)
}

func (this logger) Warnln(args ...interface{}) {
	this.WarnlnWithFields(func() *LogFields {
		return NewLogFields()
	}, args)
}

func (this logger) WarnlnWithFields(fieldsFunc FieldsFunc, args ...interface{}) {
	entry := this.WithFields(fieldsFunc().Get())
	entry.Data["file"] = fileInfo(3)
	entry.Warnln(args...)
}

func (this logger) Fatal(args ...interface{}) {
	this.FatalWithFields(func() *LogFields {
		return NewLogFields()
	}, args)
}

func (this logger) FatalWithFields(fieldsFunc FieldsFunc, args ...interface{}) {
	entry := this.WithFields(fieldsFunc().Get())
	entry.Data["file"] = fileInfo(3)
	entry.Fatal(args...)
}

type FieldsFunc func() *LogFields

func NewLogFields() *LogFields {
	fields := make(logrus.Fields)
	fields["service"] = Conf.ServiceName
	return &LogFields{
		fields: fields,
	}
}

type LogFields struct {
	fields logrus.Fields
}

func (this LogFields) Get() logrus.Fields {
	return this.fields
}

func (this LogFields) Append(k string, v interface{}) *LogFields {
	this.fields[k] = v
	return &this
}

func newLogger() logger {

	l := logger{
		logrus.New(),
	}
	l.Formatter = logFormatter

	logFile, err := os.OpenFile(
		LogPath+"/app.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		os.ModePerm,
	)

	if err != nil {
		logrus.Fatal("log file create failed.", err)
	}

	l.Out = logFile
	//l.Out = os.Stdout

	return l
}

func fileInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = 0
	}
	return fmt.Sprintf("%s:%d", file, line)
}
