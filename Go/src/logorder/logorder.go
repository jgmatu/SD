package main

import (
	"logiclog"
	"os"
)

func main() {
	output := true // Normal program yes output to stdout.
	logiclog.Order("order.txt", os.Args[1:] , output)
}
