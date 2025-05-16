// Package client provides HTTP clients for external enrichment APIs.
package client

import (
	"context"
	"fmt"
	"github.com/51mans0n/effective-mobile-task/internal/model"
	"net/http"
)

type agifyResp struct {
	Age int `json:"age"`
}

type agify struct{}

// NewAgify returns an instance of Agify client.
func NewAgify() Enricher { return &agify{} }

func (a *agify) Enrich(ctx context.Context, name string) (*model.Enriched, error) {
	key := "agify:" + name
	if v, ok := cacheGet(key); ok {
		return v.(*model.Enriched), nil
	}

	req, _ := http.NewRequestWithContext(ctx, "GET",
		fmt.Sprintf("https://api.agify.io/?name=%s", q(name)), nil)

	var r agifyResp
	if err := doJSON(req, &r); err != nil {
		return nil, err
	}

	out := &model.Enriched{Age: &r.Age}
	cacheSet(key, out)
	return out, nil
}
