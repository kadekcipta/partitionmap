package partitionmap

import (
	"fmt"
	"testing"
)

func createData() []struct {
	k string
	v interface{}
} {
	var data = []struct {
		k string
		v interface{}
	}{}
	for n := 0; n < 1000000; n++ {
		kv := struct {
			k string
			v interface{}
		}{
			fmt.Sprintf("Key#%d", n),
			n,
		}
		data = append(data, kv)
	}
	return data
}

var data = createData()

func doTask(b *testing.B, sg SetGetter) {
	for _, kv := range data {
		sg.Set(kv.k, kv.v)
		sg.Get(kv.k)
	}
}

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
