package main

import (
	"flag"
	"fmt"

	"knurov.ru/el/1c2el/helper"
)

func main() {
	configName := flag.String("config", "./config.yaml", "Set config file path")
	dryRun := flag.Bool("dry-run", false, "Dry run mode. Don`t persist data to DB")
	flag.Parse()
	fmt.Printf("Starting with! %v\n", *configName)

	hlp := helper.NewHelper(*configName, *dryRun)
	hlp.Log.Debug("Use config - %v", *configName)
	defer hlp.Destroy()
	readDir(hlp)
}
