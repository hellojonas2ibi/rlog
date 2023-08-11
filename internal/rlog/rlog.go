package rlog

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type rlogger struct {
	dir      string
	interval time.Duration
	entries  chan Entry
	stop     chan bool
	Exit     chan bool
}

type Entry struct {
	identifier string
	Time       string `json:"time"`
	Level      string `json:"level"`
	Message    string `json:"message"`
}

func New(dir string, buffer int, interval time.Duration) *rlogger {
	logger := rlogger{
		dir:      dir,
		entries:  make(chan Entry),
		interval: interval,
		Exit:     make(chan bool),
	}

	go logger.start()

	return &logger
}

func (l *rlogger) logFile(identifier string) (*os.File, error) {
	date := time.Now()
	year := fmt.Sprint(date.Year())
	month := fmt.Sprint(int(date.Month()))
	day := fmt.Sprint(date.Day())
	logDir := filepath.Join(l.dir, year, month, day)
	filename := filepath.Join(logDir, fmt.Sprintf("%s_%s-%s-%s.log", identifier, year, month, day))

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err == nil {
		return file, nil
	}

	stat, err := os.Stat(logDir)

	if (err != nil && os.IsNotExist(err)) || !stat.IsDir() {
		if err = os.MkdirAll(logDir, os.ModePerm); err != nil {
			return nil, err
		}
	}

	file, err = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)

	if err != nil {
		return nil, err
	}

	return file, nil
}

func (l *rlogger) Submit(entry Entry) {
	l.entries <- entry
}

func (l *rlogger) Close() {
	close(l.entries)
}

func (l *rlogger) start() {
	func() {
		ticker := time.NewTicker(l.interval)
		entries := make(map[string][]Entry)

		for {
			select {
			case <-ticker.C:
				fmt.Printf("-----------------------> Processing\n")
				batchWrite(l, entries)
				entries = make(map[string][]Entry)
			case entry := <-l.entries:
				if entry.identifier == "" {
					continue
				}
				entries[entry.identifier] = append(entries[entry.identifier], entry)
				fmt.Printf("--------> Adding\n")
			case <-l.stop:
				l.Exit <- true
				break
			}
		}
	}()
}

func batchWrite(logger *rlogger, entryGroups map[string][]Entry) {
	if len(entryGroups) == 0 {
		return
	}

	for id, entries := range entryGroups {
		file, err := logger.logFile(id)

		if err != nil {
			fmt.Printf("[ERROR] error opening file with id '%s'", id)
			continue
		}

		defer file.Close()

		for _, e := range entries {
			fmt.Fprintln(file, e.Message)
		}
		delete(entryGroups, id)
	}
}
