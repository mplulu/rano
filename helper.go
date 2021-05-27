package rano

import (
	"fmt"
	"runtime"
)

func GetStack() string {
	trace := make([]byte, 8192)
	count := runtime.Stack(trace, false)
	content := fmt.Sprintf("Dump (%d bytes):\n %s \n", count, trace[:count])
	return content
}
