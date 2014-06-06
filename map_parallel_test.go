// +build parallel

package partitionmap

import (
	"testing"
)

func BenchmarkPartitionMap(b *testing.B) {
	pm := New(128)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			doTask(b, pm)
		}
	})
}

func BenchmarkNonPartitionMap(b *testing.B) {
	m := &partition{m: map[string]interface{}{}}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			doTask(b, m)
		}
	})
}
