package main

import (
	"flag"
	"strings"
	"tezosign/api"
	"tezosign/common/log"
	"tezosign/common/modules"
	"tezosign/conf"
	"tezosign/infrustructure"
	"tezosign/services"

	"github.com/roylee0704/gron"
	"go.uber.org/zap"

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

	cron := gron.New()

	for k := range cfg.Networks {
		networkContext, err := provider.GetNetworkContext(cfg.Networks[k].Name)
		if err != nil {
			log.Fatal("Cron init: ", zap.Error(err))
		}

		services.AddToCron(cron, cfg, networkContext, cfg.Networks[k].Name)
	}

	cron.Start()
	defer cron.Stop()

	a := api.NewAPI(cfg, provider)
	mds := []modules.Module{a}

	modules.Run(mds)

	var gracefulStop = make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM, syscall.SIGINT)

	<-gracefulStop
	modules.Stop(mds)
}
