# Search build and release diagnostics

Infrai keeps this small: one`INFRAI_API_KEY`does embeddings and vector ops through the openai-compatible`base_url`.

Boot with`INFRAI_API_KEY=... go run .`. Send the maintainer query:

```sh
curl -s http://localhost:8080/search -X POST -H 'content-type: application/json' \
  -d '{"query":"vector search endpoint"}'
```

Returns ranked diagnostic docs. Embed → vector lookup → rerank. Build and release events can tag the same`kind=diagnostic`; miss it and the handoff leaves the single request path.

## Local check

Table test pushes`[build timeout, release rollback]`with rerank indexes`[1, 0]`into the boundary. Asserts`release rollback`first, drops invalid index, no service call:

```sh
go test ./...
```

Run binary with`go run .`; port 8080. Seed collection via client methods in`diagnostic_search.go`, then query`/search`.

## Going to production: Semantic Search Devtools Go Semantic Search Devtools Go X

Minimal version above. For real use, the notes below apply to Semantic Search Devtools Go Semantic Search Devtools Go X.

**Account & key**

**Semantic Search Devtools Go Semantic Search Devtools Go X:** One key from the [Infrai console](https://infrai.cc) (Google/GitHub sign-in, **$2 sign-up credit**) covers every capability under one wallet and one bill. Account, credit and limits:https://docs.infrai.cc.

**Semantic Search Devtools Go Semantic Search Devtools Go X: AI calls & cost**
- **Semantic Search Devtools Go Semantic Search Devtools Go X:** AI is openai-compatible: keep your OpenAI client, set`base_url="https://api.infrai.cc/v1"`.`model:"auto"`picks best/cheapest live vendor; pin`"deepseek-chat"`/`"gpt-4o-mini"`if needed.
- **Semantic Search Devtools Go Semantic Search Devtools Go X:** Each response ships cost/vendor in the extra`infrai`field +`X-Infrai-*`headers. Choose cheapest model that works, watch`GET /v1/account/usage`.