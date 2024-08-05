// Package logger exports preferred vtpl logger
package logger

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LogConfig holds the configurations for logger
type LogConfig struct {
	AppName       string `json:"appName,omitempty" yaml:"appName,omitempty"`
	SessionFolder string `json:"sessionFolder,omitempty" yaml:"sessionFolder,omitempty"`
	ConsoleLog    bool   `json:"consoleLog,omitempty" yaml:"consoleLog,omitempty"`
}

// Logger holds a Logger object
type Logger struct {
	zerolog.Logger
	// MemoryLog to view logs from API
	MemoryLog *circularBuffer
}

// New logger is created depending upon the requirement
func New(logConfig LogConfig) (*Logger, error) {
	var logWriter io.Writer
	if logConfig.ConsoleLog {
		logWriter = &zerolog.ConsoleWriter{Out: os.Stdout}
	} else {
		if logConfig.AppName == "" {
			logConfig.AppName = "ojana"
		}
		if logConfig.SessionFolder == "" {
			logConfig.SessionFolder = "session"
		}
		logFileName := fmt.Sprintf("%s.log", logConfig.AppName)

		logFilePath := path.Clean(path.Join(logConfig.SessionFolder, logFileName))

		logWriter = &lumberjack.Logger{
			Filename:   logFilePath,
			MaxSize:    10, // megabytes
			MaxBackups: 3,  //nolint:gomnd // number of files to keep
			MaxAge:     28, //nolint:gomnd // days
		}
	}
	memoryLog := newBuffer(16)
	writer := zerolog.MultiLevelWriter(logWriter, memoryLog)
	logger := &Logger{
		zerolog.New(writer).With().Timestamp().Logger(),
		memoryLog,
	}
	return logger, nil
}

const chunkSize = 1 << 16

type circularBuffer struct {
	chunks [][]byte
	r, w   int
}

func newBuffer(chunks int) *circularBuffer {
	b := &circularBuffer{chunks: make([][]byte, 0, chunks)}
	// create first chunk
	b.chunks = append(b.chunks, make([]byte, 0, chunkSize))
	return b
}

func (b *circularBuffer) Write(p []byte) (n int, err error) {
	n = len(p)

	// check if chunk has size
	if len(b.chunks[b.w])+n > chunkSize {
		// increase write chunk index
		if b.w++; b.w == cap(b.chunks) {
			b.w = 0
		}
		// check overflow
		if b.r == b.w {
			// increase read chunk index
			if b.r++; b.r == cap(b.chunks) {
				b.r = 0
			}
		}
		// check if current chunk exists
		if b.w == len(b.chunks) {
			// allocate new chunk
			b.chunks = append(b.chunks, make([]byte, 0, chunkSize))
		} else {
			// reset len of current chunk
			b.chunks[b.w] = b.chunks[b.w][:0]
		}
	}

	b.chunks[b.w] = append(b.chunks[b.w], p...)
	return
}

func (b *circularBuffer) WriteTo(w io.Writer) (n int64, err error) {
	for i := b.r; ; {
		var nn int
		if nn, err = w.Write(b.chunks[i]); err != nil {
			return
		}
		n += int64(nn)

		if i == b.w {
			break
		}
		if i++; i == cap(b.chunks) {
			i = 0
		}
	}
	return
}

func (b *circularBuffer) Reset() {
	b.chunks[0] = b.chunks[0][:0]
	b.r = 0
	b.w = 0
}
