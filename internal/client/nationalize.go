package client

import (
	"context"
	"fmt"
	"github.com/51mans0n/effective-mobile-task/internal/model"
	"net/http"
)

type nationalizeResp struct {
	Country []struct {
		CountryID   string  `json:"country_id"`
		Probability float32 `json:"probability"`
	} `json:"country"`
}

type nationalize struct{}

// NewNationalize returns a Nationalize API client.
func NewNationalize() Enricher { return &nationalize{} }

func (n *nationalize) Enrich(ctx context.Context, name string) (*model.Enriched, error) {
	key := "nat:" + name
	if v, ok := cacheGet(key); ok {
		return v.(*model.Enriched), nil
	}

	req, _ := http.NewRequestWithContext(ctx, "GET",
		fmt.Sprintf("https://api.nationalize.io/?name=%s", q(name)), nil)

	var r nationalizeResp
	if err := doJSON(req, &r); err != nil {
		return nil, err
	}

	var out model.Enriched
	if len(r.Country) > 0 {
		out.CountryCode = &r.Country[0].CountryID
		out.Probability = &r.Country[0].Probability
	}
	cacheSet(key, &out)
	return &out, nil
}
