package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/51mans0n/effective-mobile-task/internal/model"
	"net/http"
	"net/url"
	"time"

	"github.com/patrickmn/go-cache"
)

// общая память-кэш на 24 ч, проверка каждые 60 мин
var c = cache.New(24*time.Hour, time.Hour)

// Enricher – один любой внешний сервис
type Enricher interface {
	Enrich(ctx context.Context, name string) (*model.Enriched, error)
}

// =============== вспом-helper =================

// http-клиент с таймаутом 10 с (можно один на все клиенты)
var httpClient = &http.Client{Timeout: 10 * time.Second}

// doJSON выполняет запрос, читает 200-ответ в dst.
func doJSON(req *http.Request, dst any) error {
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("remote %s: %s", req.URL.Host, resp.Status)
	}
	return json.NewDecoder(resp.Body).Decode(dst)
}

// маленькие обёртки для кэша
func cacheGet(k string) (any, bool) { return c.Get(k) }
func cacheSet(k string, v any)      { c.SetDefault(k, v) }
func q(name string) string          { return url.QueryEscape(name) }
