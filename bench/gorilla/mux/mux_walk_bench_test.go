package mux

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func handler(w http.ResponseWriter, r *http.Request) {}

var routes = []string{
	"/",
	"/env",
	"/pprof/",
	"/runtime/gc",
	"/runtime/gcfree",
}

var results []string

func BenchmarkMuxWalk(b *testing.B) {
	r := mux.NewRouter()

	for _, pattern := range routes {
		r.HandleFunc(pattern, handler)
	}

	b.ResetTimer()

	result := []string{}

	for i := 0; i < b.N; i++ {
		paths := []string{}

		err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			path, err := route.GetPathTemplate()
			if err != nil {
				return err
			}

			paths = append(paths, path)
			return nil
		})

		require.NoError(b, err)
		result = paths
	}

	results = result
}

func BenchmarkMuxWalkWithPrefixAdd(b *testing.B) {
	r := mux.NewRouter()

	for _, pattern := range routes {
		r.HandleFunc(pattern, handler)
	}

	b.ResetTimer()

	result := []string{}
	prefix := "/debug"

	for i := 0; i < b.N; i++ {
		paths := []string{}

		err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			path, err := route.GetPathTemplate()
			if err != nil {
				return err
			}

			paths = append(paths, prefix+path)
			return nil
		})

		require.NoError(b, err)
		result = paths
	}

	results = result
}

func BenchmarkMuxWalkWithPrefixSprintf(b *testing.B) {
	r := mux.NewRouter()

	for _, pattern := range routes {
		r.HandleFunc(pattern, handler)
	}

	b.ResetTimer()

	result := []string{}
	prefix := "/debug"

	for i := 0; i < b.N; i++ {
		paths := []string{}

		err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			path, err := route.GetPathTemplate()
			if err != nil {
				return err
			}

			paths = append(paths, fmt.Sprintf("%s%s", prefix, path))
			return nil
		})

		require.NoError(b, err)
		result = paths
	}

	results = result
}

func BenchmarkRoutesPureSlice(b *testing.B) {
	result := []string{}

	for i := 0; i < b.N; i++ {
		paths := []string{
			"/",
			"/env",
			"/pprof/",
			"/runtime/gc",
			"/runtime/gcfree",
		}

		result = paths
	}

	results = result
}

func BenchmarkRoutesPureSliceWithPrefixAdd(b *testing.B) {
	result := []string{}
	prefix := "/debug"

	for i := 0; i < b.N; i++ {
		paths := []string{
			prefix + "/",
			prefix + "/env",
			prefix + "/pprof/",
			prefix + "/runtime/gc",
			prefix + "/runtime/gcfree",
		}

		result = paths
	}

	results = result
}

func BenchmarkRoutesPureSliceWithPrefixSprintf(b *testing.B) {
	result := []string{}
	prefix := "/debug"

	for i := 0; i < b.N; i++ {
		paths := []string{
			fmt.Sprintf("%s/", prefix),
			fmt.Sprintf("%s/env", prefix),
			fmt.Sprintf("%s/pprof/", prefix),
			fmt.Sprintf("%s/runtime/gc", prefix),
			fmt.Sprintf("%s/runtime/gcfree", prefix),
			fmt.Sprintf("%s/", prefix),
		}

		result = paths
	}

	results = result
}
