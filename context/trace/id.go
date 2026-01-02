package trace

import (
	"strconv"
	"time"
)

func generateID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}
