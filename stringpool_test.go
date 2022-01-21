// Package stringpool provides a sync.Pool of efficient
// strings.Builder workers that may be reused as needed,
// reducing the need to instantiate and allocate new
// builders in text heavy applications.
//
// From the Go standard library:
//
// A Builder is used to efficiently build a string using
// Write methods. It minimizes memory copying.
//
// A Pool is used to cache allocated but unused items for
// later reuse, relieving pressure on the garbage collector.
// That is, it makes it easy to build efficient, thread-safe
// free lists.
//
// Go 1.10 or later is required.
package stringpool

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

type pooler interface {
	Get() *strings.Builder
	Release(sb *strings.Builder)
}

type swimmer struct{ strings.Builder }

func (t swimmer) Get() *strings.Builder {
	return new(strings.Builder)
}

func (t swimmer) Release(sb *strings.Builder) {
	t.Reset()
}

func sbNonPool() pooler {
	return &swimmer{}
}

var (
	sb         *strings.Builder = &strings.Builder{}
	NewPool                     = New()
	out        string           = "" // global string return value
	global_n   int              = 0
	global_err error            = nil
	k          byte
)

func chr(c byte) string {
	return fmt.Sprintf("%c", c)
}

const (
	// scaling factor (in powers of 2)
	defaultMaxScalingFactor = 6
	maxScalingFactor        = 10
)

type (
	Any interface{}

	GetSetter interface {
		Get(key Any) (Any, error)
		Set(key Any, value Any) error
	}

	// Args implements a map of arguments
	Args interface {
		GetSetter
	}

	args map[Any]Any

	ArgSet []Args

	Benchmark interface{}

	Benchmarks interface {
		Name() string
		Scale() int
		Count() int
		Next() Benchmark
		Setup() setupFunc
		Cleanup() cleanupFunc
	}

	benchmark struct {
		name    string   // name of benchmark test
		argSets []ArgSet // multiple runs with multiple argSets
		want    Any      // return value wanted
		wantErr bool     // is error wanted?
	}

	benchmarkSet struct {
		name          string      // name of benchmark set
		scale         int         // current scaling factor
		counter       int         // current trial counter
		runs          []benchmark // multiple runs if multiple Arg sets
		scalingFactor int         // max scaling factor for benchmark set (1-10)
		setup         setupFunc   // function used to setup benchmarks
		cleanup       cleanupFunc // function used to cleanup benchmarks
	}

	setupFunc   = func(set *benchmarkSet) error
	cleanupFunc = func(set *benchmarkSet) error
)

func NewBenchmarkSet(name string, tests []benchmark, scalingFactor int, setup setupFunc, cleanup cleanupFunc) *benchmarkSet {
	if scalingFactor < 1 || scalingFactor > maxScalingFactor {
		scalingFactor = defaultMaxScalingFactor
	}

	if setup == nil {
		setup = defaultSetup
	}

	if cleanup == nil {
		cleanup = defaultCleanup
	}
	return &benchmarkSet{
		name:          name,
		scale:         0,
		counter:       0,
		runs:          tests,
		scalingFactor: scalingFactor,
		setup:         setup,
		cleanup:       cleanup,
	}
}

func BenchmarkAll(b *testing.B) {
	benchmarks := benchmarkSet{
		name:          "",
		scale:         0,
		counter:       0,
		runs:          []benchmark{},
		scalingFactor: defaultMaxScalingFactor,
		setup:         func(set *benchmarkSet) error { return nil },
		cleanup:       func(set *benchmarkSet) error { return nil },
	}

}

func defaultSetup(set *benchmarkSet) error {

	// set any global variables here
	// pass configuration options using default _config

	return nil
}

func defaultCleanup() error {
	return nil
}

func BenchmarkStringPool(b *testing.B) {
	benchmarks := []struct {
		name string
		pool pooler
		want string
	}{
		{"global", global, "global"},
		{"newPool", New(), "newPool"},
		{"non-pool", sbNonPool(), "non-pool"},
	}

	// set number of scaling factors and loop over them
	maxScalingFactor := defaultMaxScalingFactor
	for j := 0; j < maxScalingFactor; j++ {

		// scaling by powers of 2
		var scalingFactor = 1 << j

		// cycle through benchmark list
		for _, bb := range benchmarks {

			b.Run(bb.name+"("+strconv.Itoa(scalingFactor)+")", func(b *testing.B) {

				// setup and config
				sb = bb.pool.Get()

				// repeat benchmark b.N iterations
				for i := 0; i < b.N; i++ {

					// repeat various benchmark options
					for k = 0; k < 255; k++ {

						// scale internal repeats
						for l := 0; l < scalingFactor; l++ {

							// call to main benchmark function being tested
							_ = sb.WriteByte(k)

						}

					}

					// save to global variable to avoid compiler optimizations
					out = sb.String()

				}

				// cleanups and resets
				bb.pool.Release(sb)
			})
		}
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *StringPool
	}{
		// TODO: Add test cases.
		{"fake", New()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); got == tt.want {
				// do not use DeepEqual
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGet(t *testing.T) {
	fake := New()
	fakeGet := fake.Get()
	tests := []struct {
		name string
		want *strings.Builder
	}{
		// TODO: Add test cases.
		{"global", global.Get()},
		{"inline", fakeGet},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Get(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRelease(t *testing.T) {
	fake := New()
	type args struct {
		b *strings.Builder
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{"global", args{global.Get()}},
		{"inline", args{fake.Get()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Release(tt.args.b)
		})
	}
}
