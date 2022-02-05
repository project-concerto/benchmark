package main

import (
	"context"
	"flag"
	"fmt"
	"gap-lock/spree-checkout"
	"gap-lock/utils"
	"log"
	"math/rand"
	"os"
	"sort"
	"sync"

	"go.uber.org/zap"
)

func benchmark(makeFactory func() utils.APIFactory) {
	// Get THREADS setting.
	threads := utils.GetThreads()

	// Get TIMEOUT setting.
	// timeout := utils.GetTimeout()

	// Get result csv writer
	cw := utils.GetResultCSVWriter([]string{})

	for _, thread := range threads {
		// m := utils.NewMonitor(checkout.NewCheckoutSequentialFactory())
		// m := utils.NewMonitor(checkout.NewCheckoutBothTokenFactory())
		f := makeFactory()
		m := utils.NewMonitor(f)
		// result := m.BenchmarkTimeBased(thread, timeout_sec)
		result := m.BenchmarkRequestBased(thread)
		result.Write(cw, []string{})
	}
}

func prepare_sequential() {
	logger := zap.L().Sugar()

	keys := make([]int32, 0)
	for k := range utils.APITokens {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	f, err := os.Create("spree_order_number")
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}

	for _, idx := range keys {
		c := checkout.NewCheckout(utils.APITokens[idx])
		logger.Infow("Prepare Order", "idx", idx)
		r := c.RunPrepare(context.Background())
		if len(r.FailLatencies) != 0 {
			logger.Fatal("Failed during prepare")
		}
		f.WriteString(fmt.Sprintf("%v %v\n", idx, c.OrderNumber))
	}
}

// Prepare in a way that ensure order.id correspond to user.id
func prepare_ensure_order(add_payment bool) {
	logger := zap.L().Sugar()

	type Input struct {
		Idx int32
		c   *checkout.Checkout
	}
	type Result struct {
		Idx         int32
		OrderNumber string
	}

	inputChan := make(chan Input, 100)
	resultChan := make(chan Result, 100)

	keys := make([]int32, 0)
	for k := range utils.APITokens {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	// Create Order with one thread.
	go func() {
		for _, idx := range keys {
			api_token := utils.APITokens[idx]
			c := checkout.NewCheckout(api_token)
			logger.Infow("Create Order", "idx", idx)
			r := c.V1CreateOrder()
			if !r {
				logger.Fatal("Failed during create order")
			}
			inputChan <- Input{
				Idx: idx,
				c:   c,
			}
		}
	}()

	// Do the rest in parallel
	for i := 0; i < 16; i++ {
		go func() {
			for input := range inputChan {
				idx := input.Idx
				c := input.c

				logger.Infow("Prepare Order", "idx", idx)
				c.V1AddItem()
				c.V1Advance()
				c.V1AddAddress()
				if add_payment {
					c.V1AddPayment()
				}
				// c.V1Advance()
				resultChan <- Result{
					Idx:         input.Idx,
					OrderNumber: c.OrderNumber,
				}
			}
		}()
	}

	f, err := os.Create("spree_order_number")
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}

	// Collect results
	for i := 0; i < len(utils.APITokens); i++ {
		r := <-resultChan
		f.WriteString(fmt.Sprintf("%v %v\n", r.Idx, r.OrderNumber))
	}

	close(inputChan)
}

func prepare_payment() {
	l := zap.L()

	type task struct {
		idx       int
		api_token string
		order_num string
	}

	taskChan := make(chan task)
	stopChan := make(chan bool)

	go func() {
		for i := 0; i <= len(utils.APITokens); i++ {
			api_token := utils.APITokens[int32(i)]
			order_number := utils.OrderNumbers[int32(i)]
			taskChan <- task{
				idx:       i,
				api_token: api_token,
				order_num: order_number,
			}
		}
		close(stopChan)
	}()

	wg := sync.WaitGroup{}
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			for {
				select {
				case t := <-taskChan:
					if t.idx%100 == 0 {
						l.Info("Add Payment",
							zap.Int("idx", t.idx),
							zap.String("api_token", t.api_token),
							zap.String("order_num", t.order_num))
					}
					c := checkout.NewCheckout(t.api_token)
					c.OrderNumber = t.order_num
					c.V1AddPayment()
				case <-stopChan:
					wg.Done()
					return
				}
			}
		}()
	}

	wg.Wait()
}

func prepare_parallel() {
	logger := zap.L().Sugar()

	type Input struct {
		Idx      int32
		APIToken string
	}
	type Result struct {
		Idx         int32
		OrderNumber string
	}

	inputChan := make(chan Input, 100)
	resultChan := make(chan Result, 100)

	// Push input
	go func() {
		for idx, api_token := range utils.APITokens {
			inputChan <- Input{
				Idx:      idx,
				APIToken: api_token,
			}
		}
	}()

	// worker
	for i := 0; i < 16; i++ {
		go func() {
			for input := range inputChan {
				c := checkout.NewCheckout(input.APIToken)
				logger.Infow("Prepare Order", "idx", input.Idx)
				r := c.RunPrepare(context.Background())
				if len(r.FailLatencies) != 0 {
					logger.Fatal("Failed during prepare")
				}
				resultChan <- Result{
					Idx:         input.Idx,
					OrderNumber: c.OrderNumber,
				}
			}
		}()
	}

	f, err := os.Create("spree_order_number")
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}

	// Collect results
	for i := 0; i < len(utils.APITokens); i++ {
		r := <-resultChan
		f.WriteString(fmt.Sprintf("%v %v\n", r.Idx, r.OrderNumber))
	}

	close(inputChan)
}

func once() {
	// c := checkout.NewCheckoutBoth(utils.APITokens[9999], utils.OrderNumbers[9999])
	// c.RunV1AddAddress()
	// c.RunV1(context.Background())
	// c.RunPrepare(context.Background())
	// c.RunV1AddPayment(context.Background())

	api_token := utils.APITokens[int32(rand.Intn(len(utils.APITokens)))]
	c := checkout.NewCheckout(api_token)
	c.V1CreateOrder()
	c.V1AddItem()
	c.V1Advance()
	c.V1AddAddress()
	c.V1AddPayment()
	c.V1AddPayment()
}

func main() {
	commandPtr := flag.String("command", "benchmark", "benchmark/prepare")
	flag.Parse()

	utils.InitLogging()

	logger := zap.L().Sugar()
	if *commandPtr == "benchmark" {
		logger.Info("Running benchmark")
		benchmark(func() utils.APIFactory { return checkout.NewCheckoutBothTokenFactory() })
	} else if *commandPtr == "benchmark-no-contention" {
		logger.Info("Running benchmark without contention")
		benchmark(func() utils.APIFactory { return checkout.NewCheckoutNoContentionFactory() })
	} else if *commandPtr == "prepare" {
		logger.Info("Running prepare")
		// prepare_parallel()
		// prepare_sequential()
		prepare_ensure_order(false)
	} else if *commandPtr == "prepare-payment" {
		logger.Info("Running prepare payment")
		prepare_payment()
	} else if *commandPtr == "once" {
		logger.Info("Running once")
		once()
	} else {
		logger.Fatal("Unknown command")
	}
}
