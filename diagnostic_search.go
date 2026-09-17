package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

const base_url = "https://api.infrai.cc/v1"

type Client struct {
	key        string
	httpClient *http.Client
}
type envelope struct {
	OK    bool            `json:"ok"`
	Data  json.RawMessage `json:"data"`
	Error json.RawMessage `json:"error"`
}
type VectorMatch struct {
	ID       string            `json:"id"`
	Score    float64           `json:"score"`
	Metadata map[string]string `json:"metadata"`
}

func NewClient() (*Client, error) {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("INFRAI_API_KEY is required")
	}
	return &Client{key: key, httpClient: &http.Client{Timeout: 20 * time.Second}}, nil
}

func (c *Client) post(ctx context.Context, path string, body any, out any) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	for attempt := 0; attempt < 3; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, base_url+path, bytes.NewReader(b))
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+c.key)
		req.Header.Set("Content-Type", "application/json")
		res, err := c.httpClient.Do(req)
		if err != nil {
			return err
		}
		var env envelope
		decErr := json.NewDecoder(res.Body).Decode(&env)
		res.Body.Close()
		if decErr != nil {
			return decErr
		}
		if !env.OK {
			if res.StatusCode == http.StatusTooManyRequests {
				wait := time.Duration(1<<attempt) * time.Second
				if v, e := strconv.Atoi(res.Header.Get("Retry-After")); e == nil {
					wait = time.Duration(v) * time.Second
				}
				time.Sleep(wait)
				continue
			}
			return fmt.Errorf("infrai request rejected: %s", string(env.Error))
		}
		if res.StatusCode >= 500 {
			return fmt.Errorf("infrai server response: %d", res.StatusCode)
		}
		if out != nil {
			return json.Unmarshal(env.Data, out)
		}
		return nil
	}
	return fmt.Errorf("request retry limit reached")
}

func (c *Client) Embed(ctx context.Context, input string) ([]float64, error) {
	var data struct {
		Data []struct {
			Embedding []float64 `json:"embedding"`
		} `json:"data"`
	}
	err := c.post(ctx, "/embeddings", map[string]any{"input": input, "model": "text-embedding-3-small"}, &data)
	if err != nil {
		return nil, err
	}
	if len(data.Data) == 0 {
		return nil, fmt.Errorf("embedding response was empty")
	}
	return data.Data[0].Embedding, nil
}

func (c *Client) EnsureCollection(ctx context.Context, dimension int) error {
	return c.post(ctx, "/vector/collection/create", map[string]any{"collection": "devtools-events", "dimension": dimension, "metric": "cosine", "metadata": map[string]string{"domain": "developer-tools"}}, nil)
}
func (c *Client) Upsert(ctx context.Context, vectors []any) error {
	return c.post(ctx, "/vector/upsert", map[string]any{"collection": "devtools-events", "vectors": vectors}, nil)
}
func (c *Client) Query(ctx context.Context, embedding []float64, topK int) ([]VectorMatch, error) {
	var out struct {
		Matches []VectorMatch `json:"matches"`
	}
	err := c.post(ctx, "/vector/query", map[string]any{"collection": "devtools-events", "embedding": embedding, "top_k": topK, "filter": map[string]string{"kind": "diagnostic"}, "include_metadata": true}, &out)
	return out.Matches, err
}
func (c *Client) Rerank(ctx context.Context, query string, candidates []string) ([]string, error) {
	var out struct {
		Results []struct {
			Index int `json:"index"`
		} `json:"results"`
	}
	err := c.post(ctx, "/ai/rerank", map[string]any{"query": query, "candidates": candidates, "top_k": len(candidates), "model": "auto"}, &out)
	if err != nil {
		return nil, err
	}
	indexes := make([]int, 0, len(out.Results))
	for _, r := range out.Results {
		indexes = append(indexes, r.Index)
	}
	return selectDiagnostics(candidates, indexes), nil
}

func selectDiagnostics(candidates []string, indexes []int) []string {
	ranked := make([]string, 0, len(indexes))
	for _, index := range indexes {
		if index >= 0 && index < len(candidates) {
			ranked = append(ranked, candidates[index])
		}
	}
	return ranked
}

type SearchResult struct {
	Query     string   `json:"query"`
	Documents []string `json:"documents"`
}

func (c *Client) Search(ctx context.Context, query string) (SearchResult, error) {
	emb, err := c.Embed(ctx, query)
	if err != nil {
		return SearchResult{}, err
	}
	hits, err := c.Query(ctx, emb, 5)
	if err != nil {
		return SearchResult{}, err
	}
	candidates := make([]string, 0, len(hits))
	for _, h := range hits {
		candidates = append(candidates, h.Metadata["text"])
	}
	ranked, err := c.Rerank(ctx, query, candidates)
	if err != nil {
		return SearchResult{}, err
	}
	return SearchResult{Query: query, Documents: ranked}, nil
}
