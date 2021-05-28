/**
* Created by GoLand.
* User: link1st
* Date: 2019-08-21
* Time: 15:40
 */

package heper

import (
	"time"
)

// 时间差，纳秒
func DiffNano(startTime time.Time) (diff int64) {
	diff = int64(time.Since(startTime))
	return
}
