package client

import (
	"context"
	"encoding/json"
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

func NewNationalize() Enricher { return &nationalize{} }

func (n *nationalize) Enrich(ctx context.Context, name string) (*model.Enriched, error) {
	key := "nat:" + name
	if v, ok := cacheGet(key); ok {
		return v.(*model.Enriched), nil
	}

	req, _ := http.NewRequestWithContext(ctx, "GET",
		fmt.Sprintf("https://api.nationalize.io/?name=%s", q(name)), nil)

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("NATIONALIZE HTTP error:", err)
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Println("NATIONALIZE HTTP status:", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("nationalize: bad status %s", resp.Status)
	}

	var r nationalizeResp
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		fmt.Println("NATIONALIZE decode error:", err)
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
