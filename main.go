package main

import (
	"fmt"
	"os"

	"github.com/ljvmiranda921/burnout-barometer/cmd"
)

func main() {
	if err := cmd.NewCommand().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
