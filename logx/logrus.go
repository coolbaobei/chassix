package logx

import (
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/emicklei/go-restful/v3"
	log "github.com/sirupsen/logrus"
)

var formatter = &nested.Formatter{
	HideKeys:        true,
	FieldsOrder:     []string{"category", "component"},
	TimestampFormat: "2006-01-02 15:04:05",
}

//Logger custom logger
type Logger struct {
	*log.Logger
}

//Entry log entry
type Entry struct {
	// *Logger
	*log.Entry
}

type Config interface {
	Level() int
	ReportCaller() bool
	NoColors() bool
	CallerFirst() bool
	HideKeys() bool
	FieldsOrder() []string
	NoUppercaseLevel() bool
}

var config Config

//SetConfig set config
func SetConfig(cfg Config) {
	config = cfg
}

//New new logger
func New() *Logger {
	nLog := log.New()
	lg := &Logger{
		Logger: nLog,
	}
	if config != nil {
		level := log.Level(config.Level())
		if level >= log.DebugLevel {
			formatter.HideKeys = false
		}
		nLog.SetLevel(level)
		nLog.SetReportCaller(config.ReportCaller())

		formatter.NoColors = config.NoColors()
		formatter.HideKeys = config.HideKeys()
		if config.ReportCaller() {
			formatter.CallerFirst = config.CallerFirst()
		}
		if len(config.FieldsOrder()) > 0 {
			formatter.FieldsOrder = config.FieldsOrder()
		}
		formatter.NoUppercaseLevel = config.NoUppercaseLevel()
	}

	nLog.SetFormatter(formatter)
	return lg
}

const (
	fieldSvcKey = "service"
	fieldComKey = "component"
	fieldCatKey = "category"
)

//Service setting svc name
func (l *Logger) Service(svc string) *Entry {
	return &Entry{l.WithField(fieldSvcKey, svc)}
}

//Component setting svc name
func (l *Logger) Component(com string) *Entry {
	return &Entry{l.WithField(fieldComKey, com)}
}

//Category setting svc name
func (l *Logger) Category(name string) *Entry {
	return &Entry{l.WithField(fieldCatKey, name)}
}

//Component setting svc name
func (e *Entry) Component(com string) *Entry {
	return &Entry{e.WithField(fieldComKey, com)}
}

//Category setting svc name
func (e *Entry) Category(name string) *Entry {
	return &Entry{e.WithField(fieldCatKey, name)}
}

var defaultLogger = &Logger{
	log.StandardLogger(),
}

//StdLogger std logger
func StdLogger() *Logger {
	return defaultLogger
}
//SetReqInfo setting svc name
func (e *Entry) SetReqInfo(req *restful.Request) *Entry {
	return &Entry{e.WithFields(log.Fields{
		"req":   req.Request.Method + " " + req.Request.URL.String(),
		"reqId": req.Attribute("reqId"),
	})}
}