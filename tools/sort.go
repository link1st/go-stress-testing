// Package tools tool 包
package tools

// Uint64List 排序
type Uint64List []uint64

func (my64 Uint64List) Len() int           { return len(my64) }
func (my64 Uint64List) Swap(i, j int)      { my64[i], my64[j] = my64[j], my64[i] }
func (my64 Uint64List) Less(i, j int) bool { return my64[i] < my64[j] }
