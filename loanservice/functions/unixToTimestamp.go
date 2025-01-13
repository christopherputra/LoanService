package functions

import (
	"time"
)

var timeFormat = time.RFC3339


func UnixToTimeStamp(unix int64) (string) {
	return time.Unix(unix, 0).Format(timeFormat)

}