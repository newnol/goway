# GoWay

GoWay is a Go gateway for provider accounts behind LiteLLM. This starter serves a single mock model, `goway-demo`, for checking the HTTP and SSE contract before connecting a real provider.

```sh
export GOWAY_API_KEY="$(openssl rand -hex 32)"
go run ./cmd/gateway
```

In another terminal, set `GOWAY_API_KEY` to the same value and try the gateway:

```sh
curl -fsS http://127.0.0.1:8080/healthz
curl -fsS http://127.0.0.1:8080/v1/models \
  -H "Authorization: Bearer $GOWAY_API_KEY"
curl -fsSN http://127.0.0.1:8080/v1/chat/completions \
  -H "Authorization: Bearer $GOWAY_API_KEY" \
  -H 'Content-Type: application/json' \
  -d '{"model":"goway-demo","messages":[{"role":"user","content":"Hi"}],"stream":true}'
```

Set `stream` to `false` for a single JSON response. The mock always returns a fixed reply. `GOWAY_ADDR` overrides the default listen address `127.0.0.1:8080`.

## LiteLLM

Load [the example config](examples/litellm.yaml) in LiteLLM with `GOWAY_API_BASE` pointing to the gateway's reachable `/v1` URL and `GOWAY_API_KEY` set to the gateway key. Use the `goway-demo` alias to test the new route.

## Development

Use Go 1.26 or later.

```sh
go test -race ./...
go vet ./...
```

The first real provider will add OAuth login and refresh. Later milestones cover account eligibility, quota, and cooldown.
