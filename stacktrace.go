package vkc

import (
	"fmt"
	"runtime"
)

// Вывод трассировки стека. Используется в обработчике ошибок panic.
func Stacktrace(r any) {
	fmt.Printf("tracing: %v\n\n", r)
	buf := make([]byte, 1<<16)
	n := runtime.Stack(buf, false)
	fmt.Printf("%s\n", buf[:n])
}
