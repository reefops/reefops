package openfga

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/reefops/reefops/services/authorizer/internal/authorization"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Client struct {
	endpoint, storeID, modelID, token string
	http                              *http.Client
}

func New(endpoint, storeID, modelID, token string, timeout time.Duration) (*Client, error) {
	u, err := url.Parse(endpoint)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return nil, errors.New("invalid OpenFGA API URL")
	}
	if storeID == "" || modelID == "" || token == "" {
		return nil, errors.New("OpenFGA store, model and token are required")
	}
	if timeout <= 0 {
		return nil, errors.New("OpenFGA timeout must be positive")
	}
	return &Client{endpoint: strings.TrimRight(endpoint, "/"), storeID: storeID, modelID: modelID, token: token, http: &http.Client{Timeout: timeout, Transport: otelhttp.NewTransport(http.DefaultTransport)}}, nil
}

type tupleKey struct {
	User     string `json:"user"`
	Relation string `json:"relation"`
	Object   string `json:"object"`
}
type checkRequest struct {
	AuthorizationModelID string   `json:"authorization_model_id"`
	TupleKey             tupleKey `json:"tuple_key"`
	ContextualTuples     struct {
		TupleKeys []tupleKey `json:"tuple_keys"`
	} `json:"contextual_tuples"`
}
type checkResponse struct {
	Allowed *bool `json:"allowed"`
}

func (c *Client) Ready(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.endpoint+"/healthz", nil)
	if err != nil {
		return err
	}
	req.Header.Set("authorization", "Bearer "+c.token)
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("OpenFGA readiness: %w", err)
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 4096))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("OpenFGA readiness returned HTTP %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) Check(ctx context.Context, check authorization.Check) (bool, error) {
	body := checkRequest{AuthorizationModelID: c.modelID, TupleKey: tupleKey{User: check.User, Relation: check.Relation, Object: check.Object}}
	for _, t := range check.ContextualTuples {
		body.ContextualTuples.TupleKeys = append(body.ContextualTuples.TupleKeys, tupleKey{User: t.User, Relation: t.Relation, Object: t.Object})
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return false, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint+"/stores/"+url.PathEscape(c.storeID)+"/check", bytes.NewReader(payload))
	if err != nil {
		return false, err
	}
	req.Header.Set("authorization", "Bearer "+c.token)
	req.Header.Set("content-type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return false, fmt.Errorf("OpenFGA check: %w", err)
	}
	defer resp.Body.Close()
	limited := io.LimitReader(resp.Body, 1<<20)
	if resp.StatusCode != http.StatusOK {
		_, _ = io.Copy(io.Discard, limited)
		return false, fmt.Errorf("OpenFGA check returned HTTP %d", resp.StatusCode)
	}
	var out checkResponse
	decoder := json.NewDecoder(limited)
	if err := decoder.Decode(&out); err != nil {
		return false, fmt.Errorf("decode OpenFGA response: %w", err)
	}
	if out.Allowed == nil {
		return false, errors.New("OpenFGA response omitted allowed")
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return false, errors.New("OpenFGA response contains trailing data")
	}
	return *out.Allowed, nil
}
