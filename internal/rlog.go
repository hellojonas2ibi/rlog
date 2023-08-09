package rlog

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type rlogger struct {
	date string
	dir  string
}

func New(dir string) *rlogger {
	return &rlogger{
		dir: dir,
	}
}

func (l *rlogger) logFile(identifier string) (*os.File, error) {
	date := time.Now()
	year := fmt.Sprint(date.Year())
	month := fmt.Sprint(int(date.Month()))
	day := fmt.Sprint(date.Day())
	logDir := filepath.Join(l.dir, year, month, day)
	filename := filepath.Join(logDir, fmt.Sprintf("%s_%s-%s-%s.log", identifier, year, month, day))
	pathDate := date.Format("2006-01-02")

	if pathDate == l.date {
		file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE, os.ModePerm)
		if err != nil {
			return nil, err
		}
		return file, nil
	}

	stat, err := os.Stat(logDir)

	if (err != nil && os.IsNotExist(err)) || !stat.IsDir() {
		if err = os.MkdirAll(logDir, os.ModePerm); err != nil {
			return nil, err
		}
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE, os.ModePerm)

	if err != nil {
		return nil, err
	}

	l.date = pathDate
	return file, nil
}
