package context

import (
	"context"
	"fmt"
	"io"
	"os"
	"testing"
	"time"
)

var (
	ctx     context.Context
	elapsed time.Duration
)

var devNull io.Writer

func init() {
	devNull, _ = os.Create(os.DevNull)
}

func BenchmarkContextTimeout(b *testing.B) {
	var r context.Context

	background := context.Background()

	for i := 0; i < b.N; i++ {
		r, _ = context.WithTimeout(background, time.Second)
	}

	ctx = r
}

func BenchmarkContextTimeoutInline(b *testing.B) {
	var r context.Context

	for i := 0; i < b.N; i++ {
		r, _ = context.WithTimeout(context.Background(), time.Second)
	}

	ctx = r
}

func BenchmarkContextTimeoutTimed(b *testing.B) {
	var (
		r context.Context
		e time.Duration
	)

	background := context.Background()

	for i := 0; i < b.N; i++ {
		t := time.Now()
		r, _ = context.WithTimeout(background, time.Second)
		e = time.Since(t)
	}

	ctx = r
	elapsed = e
}

func BenchmarkContextTimeoutTimedPrintf(b *testing.B) {
	var (
		r context.Context
		e time.Duration
	)

	background := context.Background()

	for i := 0; i < b.N; i++ {
		t := time.Now()
		r, _ = context.WithTimeout(background, time.Second)
		e = time.Since(t)

		fmt.Fprintf(devNull, "iteration %v context creation time %v", i, e)
	}

	ctx = r
	elapsed = e
}

func BenchmarkContextTimeoutTimedPrintln(b *testing.B) {
	var (
		r context.Context
		e time.Duration
	)

	background := context.Background()

	for i := 0; i < b.N; i++ {
		t := time.Now()
		r, _ = context.WithTimeout(background, time.Second)
		e = time.Since(t)

		msg := fmt.Sprintf("iteration %v context creation time %v", i, e)
		fmt.Fprintln(devNull, msg)
	}

	ctx = r
	elapsed = e
}
