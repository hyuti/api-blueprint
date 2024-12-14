//go:build wireinject
// +build wireinject

package app

import "github.com/google/wire"

func initializeApp() (*App, error) {
	wire.Build(
		WithCfg,
		WithLogger,
		WithGinEngine,
		wire.Struct(
			new(App),
			"cfg",
			"logger",
			"ginEngine",
		),
	)
	return nil, nil
}
