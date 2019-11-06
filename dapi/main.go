package main

import (
	"context"
	"myproject/dapi/config"
	"myproject/dapi/initialize"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	initialize.Start(ctx, config.ReadConfig())
	initialize.Wait()
}
