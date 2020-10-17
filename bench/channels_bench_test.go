package bench

import (
	"sync"
	"testing"
)

var _ch chan bool

func BenchmarkChannelOpenClose(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ch = make(chan bool, 1)
		close(_ch)
	}
}

func BenchmarkChannelSend(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ch = make(chan bool, 1)
		_ch <- true
		close(_ch)
	}
}

func BenchmarkUnbufferedChannelSend(b *testing.B) {
	_ch = make(chan bool)

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()

		for range _ch {
		}
	}()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ch <- true
	}
	close(_ch)

	wg.Wait()
}
