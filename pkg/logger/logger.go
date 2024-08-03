package logger

import (
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

func New() *zerolog.Logger {
	log_file_path := "logger.log"
	lumberJackLogger := &lumberjack.Logger{
		Filename:   log_file_path,
		MaxSize:    10, // megabytes
		MaxBackups: 3,  //nolint:gomnd // number of files to keep
		MaxAge:     28, //nolint:gomnd // days
	}
	logger := zerolog.New(lumberJackLogger)
	return &logger
}
