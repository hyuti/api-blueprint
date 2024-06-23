package app

import (
	"fmt"
	els "github.com/elastic/go-elasticsearch/v8"
	"github.com/hyuti/api-blueprint/pkg/elasticsearch"
)

func WithEls() (*els.TypedClient, error) {
	if app.logger == nil {
		return nil, ErrLoggerEmpty
	}
	if app.cfg == nil {
		return nil, ErrCfgEmpty
	}
	es, err := elasticsearch.New(&elasticsearch.Cfg{
		Addresses: []string{app.cfg.Elastic.URL},
		Username:  app.cfg.Elastic.Username,
		Password:  app.cfg.Elastic.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot init elasticsearch: %w", err)
	}
	return es, nil
}
func Els() *els.TypedClient {
	mutex.Lock()
	defer mutex.Unlock()
	return app.elk
}
