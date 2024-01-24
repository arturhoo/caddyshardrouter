package caddyshardrouter

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/redis/go-redis/v9"
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

	claims, err := ParseJWT(tokenStr)
	if err != nil {
		return caddyhttp.Error(http.StatusUnauthorized, err)
	}

	customer, ok := claims["customer"].(string)
	if !ok {
		http.Error(w, "failed to parse customer", http.StatusBadRequest)
		return next.ServeHTTP(w, r)
	}
	r.Header.Set("X-Customer", customer)

	shard, err := rdb.Get(ctx, customer).Result()
	if err == redis.Nil {
		http.Error(w, "customer not found", http.StatusNotFound)
		return next.ServeHTTP(w, r)
	} else if err != nil {
		return caddyhttp.Error(http.StatusInternalServerError, fmt.Errorf("failed to query redis"))
	}
	caddyhttp.SetVar(r.Context(), "shard.upstream", shard)

	return next.ServeHTTP(w, r)
}

func jwtParseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var m JWTShardRouter
	return m, nil
}
