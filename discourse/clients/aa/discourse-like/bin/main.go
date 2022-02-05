package main

import (
	"context"
	"flag"
	"associated-access/discourse-like"
	"associated-access/utils"
	"strconv"

	"go.uber.org/zap"
)

func benchmark(mode string) {
	l := zap.L()

	// Get THREADS setting.
	threads := utils.GetThreads()

	// Get TIMEOUT setting.
	// timeout := utils.GetTimeout()

	// Get result csv writer
	cw := utils.GetResultCSVWriter([]string{"Topics"})

	// Benchmark!
	for _, thread := range threads {
		var m *utils.Monitor
		if mode == "new" {
			l.Info("Benchmarking new API", zap.Int("threads", thread))
			m = utils.NewMonitor(discourse.NewDiscourseNoContentionFactory(thread, false))
		} else if mode == "no-contention" {
			l.Info("Benchmarking non contention API", zap.Int("threads", thread))
			m = utils.NewMonitor(discourse.NewDiscourseNoContentionFactory(thread, true))
		} else {
			l.Info("Benchmarking old API", zap.Int("threads", thread))
			m = utils.NewMonitor(discourse.NewDiscourseFactory(thread))
		}
		// result := m.BenchmarkTimeBased(thread, timeout)
		result := m.BenchmarkRequestBased(thread)
		result.Write(cw, []string{strconv.Itoa(*discourse.TopicsNum)})
	}
}

func once() {
	f := discourse.NewDiscourseFactory(1)
	api := f.Make(0)
	api.Run(context.Background())
}

func main() {
	commandPtr := flag.String("command", "benchmark", "?")
	modePtr := flag.String("mode", "new", "?")
	flag.Parse()

	utils.InitLogging()
	l := zap.L()

	if *commandPtr == "benchmark" {
		benchmark(*modePtr)
	} else if *commandPtr == "prepare" {
		discourse.PreparePostsOwners()
	} else if *commandPtr == "once" {
		once()
	} else if *commandPtr == "likeall" {
		discourse.LikeEveryPostWithOneUser(555)
	} else {
		l.Fatal("Unknown command", zap.String("command", *commandPtr))
	}
}
