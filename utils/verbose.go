package utils

import (
	"fmt"

	"github.com/spf13/viper"
)

func Verbose(msg string) {
	if !viper.GetBool("verbose") {
		return
	}

	fmt.Println(msg)
}

func Verbosef(format string, v ...any) {
	if !viper.GetBool("verbose") {
		return
	}

	fmt.Printf(format+"\n", v...)
}
