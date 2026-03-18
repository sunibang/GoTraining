package bank

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/romangurevitch/go-training/pkg/client/bank/api"
	httppkg "github.com/romangurevitch/go-training/pkg/http"
)

// Client defines the interface for interacting with the Bank API.
type Client interface {
	GetToken(ctx context.Context, sub, scope string) (string, error)
	SetToken(token string)
	GetAccount(ctx context.Context, id string) (*api.AccountResponse, error)
	CreateAccount(ctx context.Context, owner string) (*api.AccountResponse, error)
	Transfer(ctx context.Context, req *api.CreateTransferRequest) (*api.TransferResponse, error)
}

func New(basePath string, httpClient *http.Client) Client {
	return &client{basePath: basePath, HTTPClient: httpClient}
}

type client struct {
	basePath   string
	HTTPClient *http.Client
	token      string // set by GetToken or SetToken, sent as Bearer on subsequent calls
}

func (c *client) SetToken(token string) {
	c.token = token
}

// GetToken issues a signed JWT for testing purposes.
func (c *client) GetToken(ctx context.Context, sub, scope string) (string, error) {
	urlPath, err := httppkg.GetURL(c.basePath, "v1/token", "")
	if err != nil {
		return "", err
	}

	payload, _ := json.Marshal(map[string]string{"sub": sub, "scope": scope})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlPath, bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}
	req.Header.Add(httppkg.HeaderApplicationJSON())

	body, err := httppkg.DoRequest(ctx, c.HTTPClient, req, http.StatusOK)
	if err != nil {
		return "", err
	}

	var res struct{ Token string }
	if err := json.Unmarshal(body, &res); err != nil {
		return "", err
	}
	c.token = res.Token
	return res.Token, nil
}

// GetAccount fetches an account by ID.
func (c *client) GetAccount(ctx context.Context, id string) (*api.AccountResponse, error) {
	start := time.Now()
	defer func() {
		slog.DebugContext(ctx, "GetAccount", slog.Duration("duration", time.Since(start)))
	}()

	if err := validator.New().Var(id, "required"); err != nil {
		return nil, err
	}

	urlPath, err := httppkg.GetURL(c.basePath, "v1/accounts/"+id, "")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)

	body, err := httppkg.DoRequest(ctx, c.HTTPClient, req, http.StatusOK)
	if err != nil {
		return nil, err
	}

	var res api.AccountResponse
	return &res, json.Unmarshal(body, &res)
}

// CreateAccount creates a new account.
func (c *client) CreateAccount(ctx context.Context, owner string) (*api.AccountResponse, error) {
	urlPath, err := httppkg.GetURL(c.basePath, "v1/accounts", "")
	if err != nil {
		return nil, err
	}

	payload, _ := json.Marshal(api.CreateAccountRequest{Owner: owner})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlPath, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Add(httppkg.HeaderApplicationJSON())
	req.Header.Set("Authorization", "Bearer "+c.token)

	body, err := httppkg.DoRequest(ctx, c.HTTPClient, req, http.StatusCreated)
	if err != nil {
		return nil, err
	}

	var res api.AccountResponse
	return &res, json.Unmarshal(body, &res)
}

// Transfer initiates a fund transfer between accounts.
func (c *client) Transfer(ctx context.Context, req *api.CreateTransferRequest) (*api.TransferResponse, error) {
	// TODO: Implement the transfer request logic.
	// Use GetAccount and CreateAccount above as a reference.
	return nil, nil
}
