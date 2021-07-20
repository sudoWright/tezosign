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
		log.Info("Sheduling operations saver every", zap.Duration("sec", dur))
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

	if conf.Cron.Assets > 0 {
		dur := time.Duration(conf.Cron.Assets) * time.Second
		log.Info("Sheduling assets saver every", zap.Duration("sec", dur))
		cron.AddFunc(gron.Every(dur), func() {

			service := New(repos.New(n.Db), repos.New(n.IndexerDB), n.Client, nil, network)

			count, err := service.AssetsIncomeOperations()
			if err != nil {
				log.Error("AssetsIncomeOperations failed", zap.Error(err))
				return
			}
			log.Info("Asset operations", zap.Uint64("count", count))
		})
	} else {
		log.Info("no sheduling assets due to missing Assets in config")
	}
}
