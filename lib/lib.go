package lib

import (
	"github.com/mrusme/xbsapi/ent"
	"go.uber.org/zap"
)

var VERSION string

type XBSContext struct {
	Config    *Config
	EntClient *ent.Client
	Logger    *zap.Logger
}
