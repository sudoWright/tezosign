package conf

import (
	log "github.com/sirupsen/logrus"
	"msig/common/baseconf"
	"msig/common/baseconf/types"
	"msig/models"
	"msig/services/rpc_client/client"
)

type (
	Config struct {
		API      API
		LogLevel string
		Networks []Network
	}

	API struct {
		ListenOnPort       uint64
		CORSAllowedOrigins []string
	}

	Network struct {
		Name    models.Network
		Params  types.DBParams
		AuthKey string
		NodeRpc client.TransportConfig
	}
)

const (
	Service         = "msig"
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
