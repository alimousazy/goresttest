package chatserver

import (
	"fmt"
	"os"
)

type Logger struct {
	file *os.File
	In   chan string
}

func (l *Logger) start() {
	for txt := range l.In {
		l.log(txt)
	}
}
func (l *Logger) log(text string) error {
	if _, err := l.file.WriteString(text); err != nil {
		fmt.Fprintf(os.Stderr, "Error writting to log file, %s.\n", err)
		return err
	}
	return nil
}

func NewLogger(path string) (*Logger, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return nil, err
	}
	return &Logger{
		file: file,
		In:   make(chan string, 1000),
	}, nil
}
