// package partitionmap provides very simple partitioned map functionalities
//
// Each partition has its own lock and any value can be distributed among partition
// based on the hashing function
package partitionmap

import (
	"hash"
	"hash/crc32"
	"sync"
)

type SetGetter interface {
	Set(string, interface{})
	Get(string) interface{}
}

type StringMap map[string]interface{}

type partition struct {
	sync.RWMutex
	m StringMap
}

func (p *partition) Set(key string, v interface{}) {
	p.Lock()
	defer p.Unlock()
	p.m[key] = v
}

func (p *partition) Get(key string) interface{} {
	p.Lock()
	defer p.Unlock()
	return p.m[key]
}

func newPartition() *partition {
	return &partition{
		m: StringMap{},
	}
}

type PartitionMap struct {
	partitions    []*partition
	hf            hash.Hash32
	partitionSize uint32
}

func (pm *PartitionMap) hash(key string) uint32 {
	pm.hf.Reset()
	pm.hf.Write([]byte(key))
	idx := pm.hf.Sum32() % pm.partitionSize
	return idx
}

func (pm *PartitionMap) Set(key string, v interface{}) {
	hk := pm.hash(key)
	pm.partitions[hk].Set(key, v)
}

func (pm *PartitionMap) Get(key string) interface{} {
	hk := pm.hash(key)
	return pm.partitions[hk].Get(key)
}

func New(partitionSize uint32, vars ...hash.Hash32) *PartitionMap {
	if partitionSize <= 0 {
		panic("Partition size must be greater 0")
	}
	hf := crc32.NewIEEE()
	if len(vars) > 0 {
		hf = vars[0]
	}
	partitions := make([]*partition, partitionSize)
	for i, _ := range partitions {
		partitions[i] = newPartition()
	}
	return &PartitionMap{
		hf:            hf,
		partitions:    partitions,
		partitionSize: partitionSize,
	}
}
