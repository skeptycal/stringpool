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
	"reflect"
	"strings"
	"sync"
	"testing"
)

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
	type args struct {
		b *strings.Builder
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{"global", args{global.Get()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Release(tt.args.b)
		})
	}
}

func TestStringPool_Release(t *testing.T) {
	type fields struct {
		pool sync.Pool
	}
	type args struct {
		b *strings.Builder
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bp := &StringPool{
				pool: tt.fields.pool,
			}
			bp.Release(tt.args.b)
		})
	}
}
