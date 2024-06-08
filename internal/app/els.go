package app

import (
	"fmt"
	els "github.com/elastic/go-elasticsearch/v8"
	"github.com/hyuti/api-blueprint/pkg/elasticsearch"
)

func WithEls() error {
	if app.logger == nil {
		return ErrLoggerEmpty
	}
	if app.cfg == nil {
		return ErrCfgEmpty
	}
	es, err := elasticsearch.New(&elasticsearch.Cfg{
		Addresses: []string{app.cfg.Elastic.URL},
		Username:  app.cfg.Elastic.Username,
		Password:  app.cfg.Elastic.Password,
	})
	if err != nil {
		return fmt.Errorf("cannot init elasticsearch: %w", err)
	}
	app.elk = es
	return nil
}
func Els() *els.TypedClient {
	mutex.Lock()
	defer mutex.Unlock()
	return app.elk
}
