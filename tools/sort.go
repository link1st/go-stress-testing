package tools

type MyUint64List []uint64

func (my64 MyUint64List) Len() int           { return len(my64) }
func (my64 MyUint64List) Swap(i, j int)      { my64[i], my64[j] = my64[j], my64[i] }
func (my64 MyUint64List) Less(i, j int) bool { return my64[i] < my64[j] }

// 判断int数据奇偶  奇数为true 偶数为false
func IsOddOrEven(number int) bool {
	if number%2 == 0 {
		return true
	}
	return false
}

// int数据取平均估值
func Average(xs []float64) (avg float64) {
	sum := 0.0
	switch len(xs) {
	case 0:
		avg = 0
	default:
		for _, v := range xs {
			sum += v
		}
		avg = sum / float64(len(xs))
	}
	return
}
