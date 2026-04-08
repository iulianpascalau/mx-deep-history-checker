package accountschecker

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type httpClient struct {
	baseURL string
	token   string
	client  *http.Client
}

// NewAPIClient creates a new API client to communicate with the network API.
func NewAPIClient(baseURL, token string) APIClient {
	return &httpClient{
		baseURL: baseURL,
		token:   token,
		client: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 100,
				IdleConnTimeout:     90 * time.Second,
			},
		},
	}
}

type statusResponse struct {
	Data struct {
		Status struct {
			ErdNonce uint64 `json:"erd_nonce"`
		} `json:"status"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}

func (c *httpClient) GetHighestNonce(ctx context.Context) (uint64, error) {
	url := fmt.Sprintf("%sv1/%s/network/status/4294967295", c.baseURL, c.token)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var res statusResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return 0, err
	}

	if res.Error != "" {
		return 0, fmt.Errorf("api error: %s", res.Error)
	}

	return res.Data.Status.ErdNonce, nil
}

type accountResponse struct {
	Data struct {
		Account AccountData `json:"account"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}

func (c *httpClient) GetAccountData(ctx context.Context, address string, nonce uint64) (*AccountData, error) {
	url := fmt.Sprintf("%sv1/%s/address/%s?blockNonce=%d", c.baseURL, c.token, address, nonce)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var res accountResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	if res.Error != "" {
		return nil, fmt.Errorf("api error: %s", res.Error)
	}

	return &res.Data.Account, nil
}
