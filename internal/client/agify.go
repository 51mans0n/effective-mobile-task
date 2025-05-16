package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/51mans0n/effective-mobile-task/internal/model"
	"net/http"
)

type agifyResp struct {
	Age int `json:"age"`
}

type agify struct{}

func NewAgify() Enricher { return &agify{} }

func (a *agify) Enrich(ctx context.Context, name string) (*model.Enriched, error) {
	key := "agify:" + name
	if v, ok := cacheGet(key); ok {
		return v.(*model.Enriched), nil
	}

	req, _ := http.NewRequestWithContext(ctx, "GET",
		fmt.Sprintf("https://api.agify.io/?name=%s", q(name)), nil)

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("AGIFY HTTP error:", err)
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Println("AGIFY HTTP status:", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("agify: bad status %s", resp.Status)
	}

	var r agifyResp
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		fmt.Println("AGIFY decode error:", err)
		return nil, err
	}

	out := &model.Enriched{Age: &r.Age}
	cacheSet(key, out)
	return out, nil
}
