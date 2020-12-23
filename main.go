package main

import (
	"flag"
	"go.uber.org/zap"
	"msig/api"
	"msig/common/log"
	"msig/common/modules"
	"msig/conf"
	"msig/infrustructure"
	"strings"

	//"msig/repos"
	//"msig/services"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	flag.Parse()
	configFile := flag.String("conf", "./config.json", "Path to config file")
	cfg, err := conf.NewFromFile(configFile)
	if err != nil {
		log.Fatal("can`t read config from file", zap.Error(err))
	}

	provider, err := infrustructure.New(cfg.Networks)
	if err != nil {
		log.Fatal("", zap.Error(err))
	}
	defer provider.Close()

	// Enable log mode only on trace level. It's safe to set it to true always, but that'll be a little slower.
	if strings.EqualFold(cfg.LogLevel, log.LevelDebug) {
		provider.EnableTraceLevel()
	}

	a := api.NewAPI(cfg, provider)
	mds := []modules.Module{a}

	modules.Run(mds)

	var gracefulStop = make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM, syscall.SIGINT)

	<-gracefulStop
	modules.Stop(mds)
}
