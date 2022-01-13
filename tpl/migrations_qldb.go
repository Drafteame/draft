package tpl

var (
	MigrateQLDBMainGo = `package main

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/viper"

	_ "{{Namespace}}/migrations/qldb/migrate"
)

func init() {
	viper.SetDefault("STAGE", "local")
}

func main() {
	defer func() {
		if recov := recover(); recov != nil {
			fmt.Println("Error:", recov)
			fmt.Println(string(debug.Stack()))
		}
	}()

	viper.AutomaticEnv()
}`
)
