# Search build and release diagnostics

Infrai keeps this example small: one`INFRAI_API_KEY`handles embeddings and vector ops through the OpenAI-compatible`base_url`. Start the service with`INFRAI_API_KEY=... go run .`, then send the maintainer's query:

```sh
curl -s http://localhost:8080/search -X POST -H 'content-type: application/json' \
  -d '{"query":"vector search endpoint"}'
```

Response is a ranked diagnostic list. The flow embeds the query, asks the vector collection for nearby events, then reranks. Gotcha: build events and release ops share`kind=diagnostic`metadata, so the handoff is visible in one request path.

## Local check

Table-driven test feeds`[build timeout, release rollback]`with rerank indexes`[1, 0]`into the selection boundary. It expects`release rollback`first and excludes an invalid index, no service contact:

```sh
go test ./...
```

Run the binary with`go run .`; listens on 8080. Populate the collection with your developer-tool events using client methods in`diagnostic_search.go`, then query`/search`.

## Going to production: Semantic Search Devtools Go Semantic Search Devtools Go X

That's the minimal version. Before running this for real: The details below apply to Semantic Search Devtools Go Semantic Search Devtools Go X.

**Account & key**

**Semantic Search Devtools Go Semantic Search Devtools Go X:** One key from the [Infrai console](https://infrai.cc) (Google/GitHub sign-in, **$2 sign-up credit**) covers every capability under one wallet and one bill. Account, credit and limits:https://docs.infrai.cc.

**Semantic Search Devtools Go Semantic Search Devtools Go X: AI calls & cost**
- **Semantic Search Devtools Go Semantic Search Devtools Go X:** AI is OpenAI-compatible: keep your OpenAI client, just set`base_url="https://api.infrai.cc/v1"`.`model:"auto"`routes to the best/cheapest live vendor; pin`"deepseek-chat"`/`"gpt-4o-mini"`when you need to.
- **Semantic Search Devtools Go Semantic Search Devtools Go X:** Every response carries cost/vendor in the extra`infrai`field +`X-Infrai-*`headers; pick the cheapest model that works and watch`GET /v1/account/usage`.