package core

import (
	"github.com/smiletrl/gateway/pkg/accesslog"
	"github.com/smiletrl/gateway/pkg/logger"
)

// core provider for other providers, such as database//jwt/logger provider.
type Provider struct {
	Access accesslog.Provider
	Logger logger.Provider
}

func BuildProvider() *Provider {
	l := logger.NewProvider()
	return &Provider{
		Access: accesslog.NewProvider(l),
		Logger: l,
	}
}
