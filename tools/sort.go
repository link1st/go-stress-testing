package tools

type MyUint64List []uint64

func (my64 MyUint64List) Len() int           { return len(my64) }
func (my64 MyUint64List) Swap(i, j int)      { my64[i], my64[j] = my64[j], my64[i] }
func (my64 MyUint64List) Less(i, j int) bool { return my64[i] < my64[j] }
