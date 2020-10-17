package bench

import (
	"fmt"
	"testing"
)

var res string

func BenchmarkStringFmt(b *testing.B) {
	var r string

	prefix := "hello"

	for i := 0; i < b.N; i++ {
		r = fmt.Sprintf("%s/world", prefix)
	}

	res = r
}

func BenchmarkStringAdd(b *testing.B) {
	var r string

	prefix := "hello"

	for i := 0; i < b.N; i++ {
		r = prefix + "/world"
	}

	res = r
}
