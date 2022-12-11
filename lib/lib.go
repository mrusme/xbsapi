package lib

import (
	"entgo.io/ent/examples/m2m2types/ent"
	"go.uber.org/zap"
)

var VERSION string

type XBSContext struct {
	Config    *Config
	EntClient *ent.Client
	Logger    *zap.Logger
}
