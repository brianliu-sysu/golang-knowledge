package main

import (
	"fmt"
	"strings"
)

func main() {
	builder := strings.Builder{}
	builder.WriteString("Hello, World!")
	builder.Grow(10)
	builder.WriteString("Hello, World!")
	builder.WriteString("Hello, World!")
	fmt.Println(builder.String())
}
