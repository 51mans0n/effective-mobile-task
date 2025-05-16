package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/51mans0n/effective-mobile-task/internal/model"
	"net/http"
)

type genderizeResp struct {
	Gender      string  `json:"gender"`
	Probability float32 `json:"probability"`
}

type genderize struct{}

func NewGenderize() Enricher { return &genderize{} }

func (g *genderize) Enrich(ctx context.Context, name string) (*model.Enriched, error) {
	key := "genderize:" + name
	if v, ok := cacheGet(key); ok {
		return v.(*model.Enriched), nil
	}

	req, _ := http.NewRequestWithContext(ctx, "GET",
		fmt.Sprintf("https://api.genderize.io/?name=%s", q(name)), nil)

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("GENDERIZE HTTP error:", err)
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Println("GENDERIZE HTTP status:", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("genderize: bad status %s", resp.Status)
	}

	var r genderizeResp
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		fmt.Println("GENDERIZE decode error:", err)
		return nil, err
	}

	out := &model.Enriched{
		Gender:      &r.Gender,
		Probability: &r.Probability,
	}
	cacheSet(key, out)
	return out, nil
}
