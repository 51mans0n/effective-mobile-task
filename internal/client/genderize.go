package client

import (
	"context"
	"fmt"
	"github.com/51mans0n/effective-mobile-task/internal/model"
	"net/http"
)

type genderizeResp struct {
	Gender      string  `json:"gender"`
	Probability float32 `json:"probability"`
}

type genderize struct{}

// NewGenderize returns a Genderize API client.
func NewGenderize() Enricher { return &genderize{} }

func (g *genderize) Enrich(ctx context.Context, name string) (*model.Enriched, error) {
	key := "genderize:" + name
	if v, ok := cacheGet(key); ok {
		return v.(*model.Enriched), nil
	}

	req, _ := http.NewRequestWithContext(ctx, "GET",
		fmt.Sprintf("https://api.genderize.io/?name=%s", q(name)), nil)

	var r genderizeResp
	if err := doJSON(req, &r); err != nil {
		return nil, err
	}

	out := &model.Enriched{
		Gender:      &r.Gender,
		Probability: &r.Probability,
	}
	cacheSet(key, out)
	return out, nil
}
