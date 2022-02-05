package utils

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

var debug = flag.Bool("debug", false, "debug mode")
var FilePrefix = flag.String("file-prefix", "result", "Result file prefix")
var Append = flag.Bool("append", false, "Append result to file")

func InitLogging() {
	// Set Logging
	var logger *zap.Logger
	var err error
	if *debug {
		logger, err = zap.NewDevelopment()
		logger.Info("Using development logger")
	} else {
		logger, err = zap.NewProduction()
		logger.Info("Using production logger")
	}
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	zap.ReplaceGlobals(logger)
}

func GetThreads() []int {
	threads := []int{1}
	threads_intput := os.Getenv("THREADS")
	if threads_intput != "" {
		threads = make([]int, 0)
		fields := strings.Split(threads_intput, ",")
		for _, field := range fields {
			thread, err := strconv.Atoi(field)
			if err != nil {
				zap.L().Fatal("Unparsable THREAD", zap.Error(err))
			}
			threads = append(threads, thread)
		}
	}
	logger := zap.L().Sugar()
	logger.Infow("Get THREADS setting", "threads", threads)

	return threads
}

func GetTimeout() int {
	timeout := 30
	timeout_intput := os.Getenv("TIMEOUT")
	if timeout_intput != "" {
		var err error
		timeout, err = strconv.Atoi(timeout_intput)
		if err != nil {
			log.Fatal(err)
		}
	}
	logger := zap.L().Sugar()
	logger.Infow("Get TIMEOUT setting", "timeout", timeout)

	return timeout
}

func GetResultCSVWriter(extra_headers []string) *csv.Writer {
	var f *os.File
	var err error
	new_file := false
	if *Append {
		filename := fmt.Sprintf("%s.csv", *FilePrefix)

		// Create file and add header if file do not exists
		if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
			zap.L().Info("Creating new CSV file")
			f, err = os.Create(filename)
			if err != nil {
				zap.L().Fatal("Can't create file", zap.Error(err))
			}
			new_file = true
		} else {
			f, err = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				zap.L().Fatal("Can't open file", zap.Error(err))
			}
		}
	} else {
		f, err = os.Create(fmt.Sprintf("%v-%v.csv", *FilePrefix, time.Now().Unix()))
		if err != nil {
			zap.L().Fatal("Can't create file", zap.Error(err))
		}
	}

	cw := csv.NewWriter(f)
	if !*Append || new_file {
		WriteResultHeader(cw, extra_headers)
	}
	return cw
}
