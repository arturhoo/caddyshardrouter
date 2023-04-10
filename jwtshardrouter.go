package caddyshardrouter

import (
	"net/http"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func init() {
	caddy.RegisterModule(JWTShardRouter{})
	httpcaddyfile.RegisterHandlerDirective("jwt_shard_router", jwtParseCaddyfile)
}

type JWTShardRouter struct {
}

func (JWTShardRouter) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.jwt_shard_router",
		New: func() caddy.Module { return new(JWTShardRouter) },
	}
}

func (m JWTShardRouter) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	authHeader := r.Header.Get("Authorization")
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	claims := ParseJWT(tokenStr)
	customer, _ := claims["customer"].(string)
	r.Header.Set("X-Customer", customer)

	shard, _ := rdb.Get(ctx, customer).Result()
	caddyhttp.SetVar(r.Context(), "shard.upstream", shard)

	return next.ServeHTTP(w, r)
}

func jwtParseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var m JWTShardRouter
	return m, nil
}
