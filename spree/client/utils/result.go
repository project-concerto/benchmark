package utils

import (
	"encoding/csv"
	"fmt"
	"sort"
)

type Result struct {
	Threads      int
	Timeout      int
	SuccessCount int
	ErrorCount   int
	Throughput   float64
	Average      int64
	P50          int64
	P90          int64
	P99          int64
	CPU          float64
}

func GetResult(results []APIResult, threads int, timeout_ms int, cpuPercent float64) Result {
	var result Result
	success_latencies := make([]int64, 0)
	fail_latencies := make([]int64, 0)
	for _, r := range results {
		if r.SuccessLatency != nil {
			success_latencies = append(success_latencies, *r.SuccessLatency)
		}
		fail_latencies = append(fail_latencies, r.FailLatencies...)
	}

	result.Threads = threads
	result.Timeout = timeout_ms / 1000
	result.SuccessCount = len(success_latencies)
	result.ErrorCount = len(fail_latencies)
	result.Throughput = float64(result.SuccessCount) * 1000 / float64(timeout_ms)
	result.Average = int64(avg(success_latencies))
	sort.Slice(success_latencies, func(i, j int) bool { return success_latencies[i] < success_latencies[j] })
	result.P50 = success_latencies[len(success_latencies)/2]
	result.P90 = success_latencies[len(success_latencies)*9/10]
	result.P99 = success_latencies[len(success_latencies)*99/100]
	result.CPU = cpuPercent
	return result
}

func WriteResultHeader(cw *csv.Writer, extra_headers []string) {
	headers := append(extra_headers,
		[]string{"Threads",
			"Timeout",
			"SuccessCount",
			"ErrorCount",
			"Throughput",
			"Average",
			"P50",
			"P90",
			"P99",
			"CPU"}...)
	cw.Write(headers)
	cw.Flush()
}

func (r *Result) Write(cw *csv.Writer, extra_values []string) {
	cw.Write(append(extra_values,
		[]string{
			fmt.Sprintf("%d", r.Threads),
			fmt.Sprintf("%d", r.Timeout),
			fmt.Sprintf("%d", r.SuccessCount),
			fmt.Sprintf("%d", r.ErrorCount),
			fmt.Sprintf("%f", r.Throughput),
			fmt.Sprintf("%d", r.Average),
			fmt.Sprintf("%d", r.P50),
			fmt.Sprintf("%d", r.P90),
			fmt.Sprintf("%d", r.P99),
			fmt.Sprintf("%f", r.CPU),
		}...))
	cw.Flush()
}

func avg(input []int64) float64 {
	total := int64(0)
	for _, v := range input {
		total += v
	}
	return float64(total) / float64(len(input))
}
