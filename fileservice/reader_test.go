package main

import (
	"testing"
)

//
// go test -bench=. -benchmem
//

// BenchmarkAsyncOneFlow-4             4443            253108 ns/op            6017 B/op         39 allocs/op
func BenchmarkAsyncOneFlow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = findFileByID(avatarGroupFiles, 123, true)
	}
}

// BenchmarkOneFlow-4                  5454            230298 ns/op            5952 B/op         35 allocs/op
func BenchmarkOneFlow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = findFileByID(avatarGroupFiles, 123, false)
	}
}

// BenchmarkAsyncListFlow-4            1100           1322728 ns/op           21991 B/op        182 allocs/op
func BenchmarkAsyncListFlow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = findListFilesByID(avatarGroupFiles, []uint64{123, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, true)
	}
}

// BenchmarkListFlow-4                  895           1452726 ns/op           20281 B/op        148 allocs/op
func BenchmarkListFlow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = findListFilesByID(avatarGroupFiles, []uint64{123, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, false)
	}
}
