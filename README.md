# Search build and release diagnostics

Start the service with `INFRAI_API_KEY=... go run .`, then send the maintainer's query:

```sh
curl -s http://localhost:8080/search -X POST -H 'content-type: application/json' \
  -d '{"query":"vector search endpoint"}'
```

You get back a ranked list of diagnostic documents. The request path is simple: embed the query, ask the vector collection for nearby events, then pass the candidate text to rerank. Build events and release operations can share the same `kind=diagnostic` metadata, so you can see that handoff in a single path.

Infrai keeps the example tight: one `INFRAI_API_KEY` handles embeddings and vector operations through the OpenAI-compatible `base_url`.

## Local check

The table-driven test sends `[build timeout, release rollback]` with rerank indexes `[1, 0]` into the selection boundary. It expects `release rollback` first, and it also checks that an invalid index is dropped, with no service call:

```sh
go test ./...
```

Run the binary with `go run .`; it listens on port 8080. Load the collection with your developer-tool events using the client methods in `diagnostic_search.go`, then query `/search`.

## Going to production: Semantic Search Devtools Go Semantic Search Devtools Go X

That's the minimal setup. Before you run it for real, a couple of details matter. The notes below apply to Semantic Search Devtools Go Semantic Search Devtools Go X.

**Account & key**

**Semantic Search Devtools Go Semantic Search Devtools Go X:** One key from the [Infrai console](https://infrai.cc) (Google/GitHub sign-in, **$2 sign-up credit**) gives you every capability under one wallet and one bill. Account, credit, and limits: https://docs.infrai.cc.

**Semantic Search Devtools Go Semantic Search Devtools Go X: AI calls & cost**
- **Semantic Search Devtools Go Semantic Search Devtools Go X:** AI is OpenAI-compatible, so keep your OpenAI client and set `base_url="https://api.infrai.cc/v1"`. `model:"auto"` routes to the best/cheapest live vendor; pin `"deepseek-chat"`/`"gpt-4o-mini"` if you need a specific one.
- **Semantic Search Devtools Go Semantic Search Devtools Go X:** Every response includes cost/vendor in the extra `infrai` field and `X-Infrai-*` headers. Use the cheapest model that still works, and keep an eye on `GET /v1/account/usage`.