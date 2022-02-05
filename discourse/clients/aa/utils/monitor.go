package utils

import (
	"context"
	"fmt"
	"log"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"go.uber.org/zap"
)

const WARMUP_REQ = 1000
const BENCHMARK_REQ = 3500
const COOLDOWN_REQ = 1000

const WARMUP_SEC = 10
const BENCHMARK_SEC = 60
const COOLDOWN_SEC = 10

type APIResult struct {
	SuccessLatency *int64
	FailLatencies  []int64
}

// API run something with retry logic.
type API interface {
	Run(ctx context.Context) APIResult
}

// This should be thread safe!
type APIFactory interface {
	Prepare()
	Make(threadId int) API
	Stop()
}

type Monitor struct {
	factory  APIFactory
	counting int32 // Workers count when `counting` is greater than 0.
}

func NewMonitor(factory APIFactory) *Monitor {
	return &Monitor{
		factory:  factory,
		counting: 0,
	}
}

func (m *Monitor) collector(resultChan chan APIResult, result chan []APIResult) {
	results := make([]APIResult, 0)
	for result := range resultChan {
		if atomic.LoadInt32(&m.counting) == 1 {
			results = append(results, result)
		}
	}
	result <- results
}

func (m *Monitor) BenchmarkTimeBased(threads, timeout_sec int) Result {
	l := zap.L().Sugar()
	l.Info("BenchmarkTimeBased")

	m.factory.Prepare()

	ctx, cancel := context.WithCancel(context.Background())
	resultChan := make(chan APIResult)
	collectorChan := make(chan []APIResult)
	var wg sync.WaitGroup

	var percent float64
	go CPUMonitor(ctx, &wg, &percent)

	go m.collector(resultChan, collectorChan)

	for i := 0; i < threads; i++ {
		go m.worker(ctx, i, resultChan, &wg)
	}

	// Warm Up, wait for WARMUP_SEC seconds.
	l.Info("Monitor Warmup")
	time.Sleep(time.Duration(WARMUP_SEC) * time.Second)
	atomic.StoreInt32(&m.counting, 1)
	// Benchmark, wait for BENCHMARK_SEC seconds.
	l.Info("Monitor Benchmark")
	time.Sleep(time.Duration(timeout_sec) * time.Second)
	atomic.StoreInt32(&m.counting, 0)
	// Cool Down, wait for COOLDOWN_SEC seconds.
	l.Info("Monitor Cooldown")
	time.Sleep(time.Duration(COOLDOWN_SEC) * time.Second)

	m.factory.Stop()
	cancel()
	wg.Wait()
	close(resultChan)

	ra := <-collectorChan
	return GetResult(ra, threads, timeout_sec*1000, percent)
}

func (m *Monitor) BenchmarkRequestBased(threads int) Result {
	l := zap.L().Sugar()
	l.Info("BenchmarkRequestBased")

	m.factory.Prepare()

	ctx, cancel := context.WithCancel(context.Background())
	resultsChan := make(chan APIResult)
	results := make([]APIResult, 500)
	var wg sync.WaitGroup

	var percent float64
	go CPUMonitor(ctx, &wg, &percent)

	for i := 0; i < threads; i++ {
		go m.worker(ctx, i, resultsChan, &wg)
	}

	collectNSuccessResult := func(n int, collectResult bool) {
		for i := 0; i < n; i++ {
			r := <-resultsChan
			if collectResult {
				results = append(results, r)
			}
			if r.SuccessLatency == nil {
				i--
			}
		}
	}

	// Warm Up, wait for WARMUP_REQ requests.
	l.Info("Monitor Warmup")
	collectNSuccessResult(WARMUP_REQ, false)

	// Benchmark, wait for BENCHMARK_REQ requests.
	start := time.Now()
	l.Info("Monitor Benchmark")
	collectNSuccessResult(BENCHMARK_REQ, true)
	benchmark_elapsed_ms := time.Since(start).Milliseconds()

	// Cool Down, wait for BENCHMARK_REQ requests.
	l.Info("Monitor CollDown")
	collectNSuccessResult(COOLDOWN_REQ, false)

	l.Info("Monitor Stopping")
	// Order matters!
	cancel()
	wg.Wait()
	m.factory.Stop()
	return GetResult(results, threads, int(benchmark_elapsed_ms), percent)
}

func (m *Monitor) Benchmark(threads int, timeout_sec int) Result {
	return m.BenchmarkRequestBased(threads)
	// return m.BenchmarkTimeBased(threads, timeout_sec)
}

func (m *Monitor) worker(ctx context.Context, id int, resultsChan chan APIResult, wg *sync.WaitGroup) {
	wg.Add(1)
	for {
		api := m.factory.Make(id)
		result := api.Run(ctx)
		select {
		case <-ctx.Done():
			wg.Done()
			return
		case resultsChan <- result:
		}
	}
}

func CPUMonitor(ctx context.Context, wg *sync.WaitGroup, cpuPercent *float64) {
	wg.Add(1)
	// Monitor CPU stats
	cputTimes1, err := cpu.Times(false)
	if err != nil {
		log.Fatal(err)
	}

	<-ctx.Done()

	cputTimes2, err := cpu.Times(false)
	if err != nil {
		log.Fatal(err)
	}

	percents, err := calculateAllBusy(cputTimes1, cputTimes2)
	if err != nil {
		log.Fatal(err)
	}

	*cpuPercent = percents[0]
	wg.Done()
}

// Copied from  github.com/shirou/gopsutil
func getAllBusy(t cpu.TimesStat) (float64, float64) {
	busy := t.User + t.System + t.Nice + t.Iowait + t.Irq +
		t.Softirq + t.Steal
	return busy + t.Idle, busy
}

// Copied from  github.com/shirou/gopsutil
func calculateBusy(t1, t2 cpu.TimesStat) float64 {
	t1All, t1Busy := getAllBusy(t1)
	t2All, t2Busy := getAllBusy(t2)

	if t2Busy <= t1Busy {
		return 0
	}
	if t2All <= t1All {
		return 100
	}
	return math.Min(100, math.Max(0, (t2Busy-t1Busy)/(t2All-t1All)*100))
}

// Copied from  github.com/shirou/gopsutil
func calculateAllBusy(t1, t2 []cpu.TimesStat) ([]float64, error) {
	// Make sure the CPU measurements have the same length.
	if len(t1) != len(t2) {
		return nil, fmt.Errorf(
			"received two CPU counts: %d != %d",
			len(t1), len(t2),
		)
	}

	ret := make([]float64, len(t1))
	for i, t := range t2 {
		ret[i] = calculateBusy(t1[i], t)
	}
	return ret, nil
}
