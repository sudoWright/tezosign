package services

import (
	"tezosign/common/log"
	"tezosign/conf"
	"tezosign/infrustructure"
	"tezosign/models"
	"tezosign/repos"
	"time"

	"go.uber.org/zap"

	"github.com/roylee0704/gron"
)

func AddToCron(cron *gron.Cron, conf conf.Config, n infrustructure.NetworkContext, network models.Network) {
	if conf.Cron.Operations > 0 {
		dur := time.Duration(conf.Cron.Operations) * time.Second
		log.Info("Sheduling counter saver every", zap.Duration("sec", dur))
		cron.AddFunc(gron.Every(dur), func() {

			service := New(repos.New(n.Db), repos.New(n.IndexerDB), n.Client, nil, network)

			count, err := service.CheckOperations()
			if err != nil {
				log.Error("CheckOperations failed", zap.Error(err))
				return
			}
			log.Info("Updated operations", zap.Int64("count", count))
		})
	} else {
		log.Info("no sheduling operations due to missing Operations in config")
	}

}
