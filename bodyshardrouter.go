package caddyshardrouter

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func init() {
	caddy.RegisterModule(BodyShardRouter{})
	httpcaddyfile.RegisterHandlerDirective("body_shard_router", bodyParseCaddyfile)
}

type BodyShardRouter struct {
}

func (BodyShardRouter) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.body_shard_router",
		New: func() caddy.Module { return new(BodyShardRouter) },
	}
}

func (m BodyShardRouter) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	// From https://github.com/caddyserver/caddy/blob/f8b59e77f83c05da87bd5e3780fb7522b863d462/modules/caddyhttp/replacer.go#L162
	if r.Body == nil {
		return next.ServeHTTP(w, r)
	}

	defer r.Body.Close()
	buf := new(bytes.Buffer)
	_, _ = io.Copy(buf, r.Body)
	r.Body = io.NopCloser(buf)

	body := buf.String()
	var data map[string]interface{}
	err := json.Unmarshal([]byte(body), &data)
	if err != nil {
		return next.ServeHTTP(w, r)
	}
	customer, ok := data["customer"].(string)
	if !ok {
		return next.ServeHTTP(w, r)
	}
	r.Header.Set("X-Customer", customer)

	shard, _ := rdb.Get(ctx, customer).Result()
	caddyhttp.SetVar(r.Context(), "shard.upstream", shard)

	return next.ServeHTTP(w, r)
}

func bodyParseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var m BodyShardRouter
	return m, nil
}
