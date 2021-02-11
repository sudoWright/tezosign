package conf

import (
	"tezosign/common/baseconf"
	"tezosign/common/baseconf/types"
	"tezosign/models"
	"tezosign/services/rpc_client/client"

	log "github.com/sirupsen/logrus"
)

type (
	Config struct {
		API      API
		LogLevel string
		Cron     Cron
		Networks []Network
	}

	API struct {
		ListenOnPort       uint64
		CORSAllowedOrigins []string
		IsProtocolHttps    bool
	}

	Cron struct {
		Operations int64
	}

	Auth struct {
		AuthKey         string
		SessionHashKey  string
		SessionBlockKey string
	}

	Network struct {
		Name          models.Network
		Params        types.DBParams
		IndexerParams types.DBParams
		Auth          Auth
		NodeRpc       client.TransportConfig
	}
)

const (
	Service         = "tezosign"
	TtlRefreshToken = 3 * 60 * 60 // 3 hours in seconds
	TtlJWT          = 1 * 60 * 60 // 1 hour in seconds
	TtlCookie       = 3 * 60 * 60 // in seconds
)

func NewFromFile(cfgFile *string) (cfg Config, err error) {
	err = baseconf.Init(&cfg, cfgFile)
	if err != nil {
		return cfg, err
	}

	err = cfg.Validate()
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

// Validate validates all Config fields.
func (config Config) Validate() error {
	//TODO add validations
	return baseconf.ValidateBaseConfigStructs(&config)
}

// DbLogger is a simple log wrapper for use with gorm and logrus.
type DbLogger struct{}

// Print directs the log ouput to trace level.
func (*DbLogger) Print(args ...interface{}) {
	log.Traceln(args...)
}
