package model

import (
	"sort"
	"strconv"
)

type TimestampManager struct {
	ObjectManager
}

func NewTimestampManager() TimestampManager {
	om := NewObjectManager()
	return TimestampManager{om}
}

func (tm *TimestampManager) Register(timestamp Timestamp) {
	tm.register(timestamp)
}

func (tm *TimestampManager) Unregister(timestamp Timestamp) {
	tm.unregister(timestamp)
}

type Uint64Slice []uint64

func (s Uint64Slice) Len() int           { return len(s) }
func (s Uint64Slice) Less(i, j int) bool { return s[i] < s[j] }
func (s Uint64Slice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

type chan_ts chan Timestamp

func (tm *TimestampManager) emitBlock(c chan_ts, block ManageableSet) {
	// Sort keys within block
	keys := make([]string, len(block))
	i := 0
	for str, _ := range block {
		keys[i] = str
	}
	sort.Strings(keys)

	// Output to chan
	for _, key := range keys {
		ts, ok := tm.GetByKey(key)
		if ok {
			c <- ts.(Timestamp)
		}
	}
}

func (tm *TimestampManager) emitTimestamps(c chan_ts, bh []uint64) {
	// Iterate through blocks
	for h := range bh {
		block := tm.GetGroup(strconv.FormatUint(uint64(h), 10))
		tm.emitBlock(c, block)
	}
}

func (tm *TimestampManager) sortedBlocks() Uint64Slice {
	blocks := tm.ObjectManager.by_group

	// Get list of block heights
	block_heights := make(Uint64Slice, len(blocks))
	i := 0
	for h := range blocks {
		int_height, err := strconv.ParseUint(h, 10, 64)
		if err != nil {
			panic(err)
		}
		block_heights[i] = int_height
		i++
	}

	// Sort and return
	sort.Sort(block_heights)
	return block_heights
}

func (tm *TimestampManager) Iter() <-chan Timestamp {
	c := make(chan Timestamp)
	go tm.emitTimestamps(c, tm.sortedBlocks())
	return c
}
