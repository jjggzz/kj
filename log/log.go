package log

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"io"
	"os"
)

// 创建一个go-kit日志对象，默认输出到控制台
// 级别默认为Debug
func BuildLogger(name string, w io.Writer, options ...level.Option) log.Logger {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		if w != nil {
			logger = log.NewLogfmtLogger(w)
		}
		logger = log.NewSyncLogger(logger)

		logger = level.NewFilter(logger, level.AllowDebug())
		if options != nil && len(options) > 0 {
			logger = level.NewFilter(logger, options...)
		}
		logger = log.With(logger,
			"svc", name,
			"ts", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}
	return logger
}
