package mymware

import (
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

func (l *Logger) Level() log.Lvl {
	switch l.Logger.Level {
	case logrus.DebugLevel:
		return log.DEBUG
	case logrus.WarnLevel:
		return log.WARN
	case logrus.ErrorLevel:
		return log.ERROR
	case logrus.InfoLevel:
		return log.INFO
	}
	return log.OFF
}

func (l *Logger) SetPrefix(s string) {
	// TODO
}

func (l *Logger) Prefix() string {
	// TODO
	return ""
}

func (l *Logger) SetHeader(h string) {
	// TODO
}

func (l *Logger) SetLevel(lvl log.Lvl) {
	switch lvl {
	case log.DEBUG:
		l.Logger.SetLevel(logrus.DebugLevel)
	case log.WARN:
		l.Logger.SetLevel(logrus.WarnLevel)
	case log.ERROR:
		l.Logger.SetLevel(logrus.ErrorLevel)
	case log.INFO:
		l.Logger.SetLevel(logrus.InfoLevel)
	default:
		l.Logger.Panic("Invalid level")
	}
}

func (l *Logger) Output() io.Writer {
	return l.Logger.Out
}

func (l *Logger) SetOutput(w io.Writer) {
	l.Logger.SetOutput(w)
}

func (l *Logger) Printj(j log.JSON) {
	l.Logger.WithFields(logrus.Fields(j)).Print()
}

func (l Logger) Debugj(j log.JSON) {
	l.Logger.WithFields(logrus.Fields(j)).Debug()
}

func (l Logger) Infoj(j log.JSON) {
	l.Logger.WithFields(logrus.Fields(j)).Info()
}

func (l Logger) Warnj(j log.JSON) {
	l.Logger.WithFields(logrus.Fields(j)).Warn()
}

func (l Logger) Errorj(j log.JSON) {
	l.Logger.WithFields(logrus.Fields(j)).Error()
}

func (l Logger) Fatalj(j log.JSON) {
	l.Logger.WithFields(logrus.Fields(j)).Fatal()
}

func (l Logger) Panicj(j log.JSON) {
	l.Logger.WithFields(logrus.Fields(j)).Panic()
}

// -------------------------------------------------
// -------------------------------------------------
// -------------------------------------------------

func logrusMiddlewareHandler(c echo.Context, next echo.HandlerFunc, skipper middleware.Skipper) error {
	if skipper != nil && skipper(c) {
		return next(c)
	}

	req := c.Request()
	res := c.Response()

	bytesIn := req.Header.Get(echo.HeaderContentLength)
	if bytesIn == "" {
		bytesIn = "0"
	}

	fields := log.JSON{
		"user_agent": req.UserAgent(),
		"remote_ip":  c.RealIP(),
		"host":       req.Host,
		// "uri":           req.RequestURI,
		"method": req.Method,
		"path":   req.URL.Path,
		"status": res.Status,
		// "referer":       req.Referer(),
		// "bytes_in":      bytesIn,
		// "bytes_out":     strconv.FormatInt(res.Size, 10),
	}
	if rid, ok := req.Context().Value(echo.HeaderXRequestID).(string); ok && rid != "" {
		fields["request_id"] = rid
	}

	c.Logger().Infoj(fields)

	if err := next(c); err != nil {
		c.Error(err)
	}
	return nil
}

// LoggerWithSkipper -
func LoggerWithSkipper(skipper middleware.Skipper) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return logrusMiddlewareHandler(c, next, skipper)
		}
	}
}
