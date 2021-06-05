/**
* Package statistics
*
* User: link1st
* Date: 2020/9/28
* Time: 14:02
 */
package statistics

import (
	"reflect"
	"testing"
)

// TestPrintMap
func TestPrintMap(t *testing.T) {

	tt := map[string]struct {
		a      map[int]int
		result string
	}{
		"test1": {a: map[int]int{
			200: 50,
			500: 10,
			100: 20,
		}, result: "100:20;200:50;500:10"},
	}

	for _, value := range tt {
		str := printMap(value.a)
		if !reflect.DeepEqual(value.result, str) {
			t.Errorf("数据不一致 预期:%v 实际:%v", value.result, str)
		}
	}
}
