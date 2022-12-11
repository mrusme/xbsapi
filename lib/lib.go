package lib

import (
	"go.uber.org/zap"
)

type XBSContext struct {
	Config *Config
	Logger *zap.Logger
}
