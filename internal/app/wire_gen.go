// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package app

// Injectors from wire.go:

func initializeApp() (*App, error) {
	config, err := WithCfg()
	if err != nil {
		return nil, err
	}
	logger, err := WithLogger()
	if err != nil {
		return nil, err
	}
	appApp := &App{
		cfg:    config,
		logger: logger,
	}
	return appApp, nil
}
