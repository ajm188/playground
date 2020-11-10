package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"sync"
	"time"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var (
	messageCh chan string
	duration  time.Duration
)

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	start := time.Now()

	ctx, cancel := context.WithTimeout(context.TODO(), duration)
	timeNow := time.Now()
	ctxElapsed := timeNow.Sub(start)
	//ctx, cancel := context.WithDeadline(ctx, deadline)
	defer cancel()

	select {
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			timeNow = time.Now()
			elapsed := timeNow.Sub(start)

			messageCh <- fmt.Sprintf("worker %v context error: %v, context diff in time: %v time creating context: %v\n", id, ctx.Err(), elapsed, ctxElapsed)
		}
	default:
		// NOOP fmt.Printf("Worker %d done\n", id)
	}
}

func main() {
	parallelism := flag.Int("parallelism", 10000, "number of workers to spawn")
	profile := flag.String("profile", "", "where to save CPU profile. no profile is taken if empty")
	flag.DurationVar(&duration, "worker-timeout", time.Millisecond, "")

	flag.Parse()

	if *profile != "" {
		f, err := os.Create(*profile)
		check(err)

		defer f.Close()

		check(pprof.StartCPUProfile(f))
		defer pprof.StopCPUProfile()
	}

	messageCh = make(chan string, *parallelism)

	start := time.Now()
	wg := sync.WaitGroup{}

	for i := 1; i <= *parallelism; i++ {
		wg.Add(1)

		go worker(i, &wg)
	}

	wg.Wait()
	close(messageCh)

	timeNow := time.Now()
	elapsed := timeNow.Sub(start)

	for msg := range messageCh {
		fmt.Print(msg)
	}

	fmt.Printf("time elapsed of the full program: %v\n", elapsed)
}
