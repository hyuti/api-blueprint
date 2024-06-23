//go:build wireinject
// +build wireinject

package app

import "github.com/google/wire"

func initializeApp() (*App, error) {
	wire.Build(
		WithCfg,
		WithLogger,
		WithTele,
		wire.Struct(
			new(App),
			"cfg",
			"logger",
			"tele",
		),
	)
	return nil, nil
}
